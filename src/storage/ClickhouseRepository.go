package storage

import (
	"context"
	"database/sql"
	"errors"
	opentracingBridge "github.com/bassbeaver/eventhouse/service/opentracing"
	"github.com/bassbeaver/logopher"
	"github.com/opentracing/opentracing-go"
	"time"
)

type clickhouseRepository struct {
	dbConnect         *sql.DB
	opentracingBridge *opentracingBridge.Bridge
}

func (cr *clickhouseRepository) Save(
	eventType string,
	idempotencyKey string,
	entityType string,
	entityId string,
	payload string,
	ctx context.Context,
) (*Event, error) {
	newEventId, newEventIdError := generateNewEventId()
	if nil != newEventIdError {
		return nil, newEventIdError
	}

	newEvent := &Event{
		EventId:    newEventId,
		EventType:  eventType,
		EntityType: entityType,
		EntityId:   entityId,
		Recorded:   time.Now(),
		Payload:    payload,
	}

	if nil != ctx {
		if parentSpan := opentracing.SpanFromContext(ctx); nil != parentSpan {
			repoSaveSpan := cr.opentracingBridge.Tracer().StartSpan(
				"event_repo__save",
				opentracing.ChildOf(parentSpan.Context()),
				opentracing.Tag{Key: "EventId", Value: newEvent.EventId},
			)
			defer repoSaveSpan.Finish()
		}
	}

	tx, txError := cr.dbConnect.Begin()
	if nil != txError {
		return nil, errors.New("Failed to create Clickhouse transaction: " + txError.Error())
	}

	// NOTICE:
	// In select statement subquery have to be at first place, because database/sql package has a bug and
	// do not recognize prepared statement placeholder in SELECT statement if it is located just after SELECT keyword
	stmt, prepareErr := tx.Prepare(
		"INSERT INTO events (PreviousEventId, EventId, EventType, IdempotencyKey, EntityType, EntityId, Recorded, Payload) " +
			"SELECT " +
			"(SELECT IF(0 < MAX(EventId), MAX(EventId), NULL) FROM events WHERE EntityType=? AND EntityId=?) as PreviousEventId, " +
			"? as EventId, " +
			"? as EventType, " +
			"? as IdempotencyKey, " +
			"? as EntityType, " +
			"? as EntityId, " +
			"? as Recorded, " +
			"? as Payload",
	)
	if nil != prepareErr {
		return nil, errors.New("Failed to prepare insert statement: " + prepareErr.Error())
	}

	_, stmtExecErr := stmt.Exec(
		newEvent.EntityType,
		newEvent.EntityId,
		newEvent.EventId,
		newEvent.EventType,
		idempotencyKey,
		newEvent.EntityType,
		newEvent.EntityId,
		newEvent.Recorded,
		newEvent.Payload,
	)
	if nil != stmtExecErr {
		return nil, errors.New("Failed to execute prepared insert statement: " + stmtExecErr.Error())
	}

	commitErr := tx.Commit()
	if nil != commitErr {
		return nil, errors.New("Failed to commit transaction: " + commitErr.Error())
	}

	return newEvent, nil
}

func (cr *clickhouseRepository) Get(eventId uint64) (*Event, error) {
	rows, queryErr := cr.dbConnect.Query(
		"SELECT EventId, PreviousEventId, EventType, EntityType, EntityId, Recorded, Payload FROM events WHERE EventId = ? LIMIT 1",
		eventId,
	)
	if nil != queryErr {
		return nil, errors.New("Failed to query Event: " + queryErr.Error())
	}

	if !rows.Next() {
		return nil, nil
	}

	event := &Event{}
	scanErr := rows.Scan(&event.EventId, &event.PreviousEventId, &event.EventType, &event.EntityType, &event.EntityId, &event.Recorded, &event.Payload)
	if nil != scanErr {
		return nil, errors.New("Failed to scan query results: " + scanErr.Error())
	}

	return event, nil
}

func (cr *clickhouseRepository) EntityStream(
	entityType string,
	entityId string,
	filterFromEventId uint64,
	includeFromEvent bool,
	loggerObj *logopher.Logger,
	ctx context.Context,
) (chan *Event, error) {
	// Buffered channel used to avoid case when stream reader is very slow and repository loaded all events from DB and pushed it to chan.
	// We want to load events with same pace as they are read.
	eventsChan := make(chan *Event, streamChanBufferSize)

	go cr.performStreamRead(
		func(lastEventId uint64) (*sql.Rows, error) {
			return cr.performEntityStreamBatchQuery(entityType, entityId, filterFromEventId, includeFromEvent, lastEventId)
		},
		eventsChan,
		loggerObj,
		ctx,
	)

	return eventsChan, nil
}

func (cr *clickhouseRepository) GlobalStream(
	filterFromEventId uint64,
	includeFromEvent bool,
	filterEntityType []string,
	filterEventType []string,
	loggerObj *logopher.Logger,
	ctx context.Context,
) (chan *Event, error) {
	eventsChan := make(chan *Event, streamChanBufferSize)

	go cr.performStreamRead(
		func(lastEventId uint64) (*sql.Rows, error) {
			return cr.performGlobalStreamBatchQuery(filterFromEventId, includeFromEvent, filterEntityType, filterEventType, lastEventId)
		},
		eventsChan,
		loggerObj,
		ctx,
	)

	return eventsChan, nil
}

func (cr *clickhouseRepository) performStreamRead(
	batchQueryPerformer func(filterFromEventId uint64) (*sql.Rows, error),
	eventsChan chan *Event,
	loggerObj *logopher.Logger,
	ctx context.Context,
) {
	var lastEventId uint64
	var rows *sql.Rows
	queryNextBatch := true

	for queryNextBatch {
		var querySpan opentracing.Span
		querySpanIsOpened := false
		closeQuerySpan := func() {
			if nil != querySpan && querySpanIsOpened {
				querySpan.Finish()
				querySpanIsOpened = false
			}
		}
		if nil != ctx {
			if parentSpan := opentracing.SpanFromContext(ctx); nil != parentSpan {
				querySpan = cr.opentracingBridge.Tracer().StartSpan(
					"repository__entity_stream_batch_query",
					opentracing.ChildOf(parentSpan.Context()),
				)
				querySpanIsOpened = true
			}
		}

		var rowsError error
		rows, rowsError = batchQueryPerformer(lastEventId)
		if nil != rowsError {
			if nil != loggerObj {
				loggerObj.Critical("Failed to execute batch query", &logopher.MessageContext{"error": rowsError.Error()})
			}
			close(eventsChan)

			return
		}
		queryNextBatch = false

		for rows.Next() {
			closeQuerySpan()

			queryNextBatch = true

			event := &Event{}
			scanErr := rows.Scan(&event.EventId, &event.PreviousEventId, &event.EventType, &event.EntityType, &event.EntityId, &event.Recorded, &event.Payload)
			if nil != scanErr {
				if nil != loggerObj {
					loggerObj.Critical("Failed to scan Event from DB response", &logopher.MessageContext{"error": scanErr.Error()})
				}
				close(eventsChan)

				return
			}

			lastEventId = event.EventId

			eventsChan <- event
		}
		if nil != rows.Err() {
			if nil != loggerObj {
				loggerObj.Critical("Failed to perform Next on queried result", &logopher.MessageContext{"error": rows.Err().Error()})
			}
			close(eventsChan)

			return
		}
	}

	close(eventsChan)
}

func (cr *clickhouseRepository) performEntityStreamBatchQuery(
	entityType,
	entityId string,
	filterFromEventId uint64,
	includeFromEvent bool,
	lastEventId uint64,
) (*sql.Rows, error) {
	sqlText :=
		`SELECT EventId, PreviousEventId, EventType, EntityType, EntityId, Recorded, Payload FROM events
		WHERE EntityType = ? AND EntityId = ? `
	sqlParams := []interface{}{entityType, entityId}

	if 0 != lastEventId {
		sqlText += " AND EventId > ? "
		sqlParams = append(sqlParams, lastEventId)
	}

	if 0 != filterFromEventId {
		if includeFromEvent {
			sqlText += " AND EventId >= ? "
		} else {
			sqlText += " AND EventId > ? "
		}
		sqlParams = append(sqlParams, filterFromEventId)
	}

	sqlText += " ORDER BY EventId ASC LIMIT " + rowsPerLoop

	rows, queryErr := cr.dbConnect.Query(sqlText, sqlParams...)
	if nil != queryErr {
		return nil, queryErr
	}

	return rows, nil
}

func (cr *clickhouseRepository) performGlobalStreamBatchQuery(
	filterFromEventId uint64,
	includeFromEvent bool,
	filterEntityType []string,
	filterEventType []string,
	lastEventId uint64,
) (*sql.Rows, error) {
	sqlText := "SELECT EventId, PreviousEventId, EventType, EntityType, EntityId, Recorded, Payload FROM events"
	whereText := ""
	sqlParams := make([]interface{}, 0)

	appendWhereOrAndToSqlText := func() {
		if "" == whereText {
			whereText += " WHERE "
		} else {
			whereText += " AND "
		}
	}

	if 0 != lastEventId {
		appendWhereOrAndToSqlText()
		whereText += " EventId > ? "
		sqlParams = append(sqlParams, lastEventId)
	}

	if 0 != filterFromEventId {
		appendWhereOrAndToSqlText()
		if includeFromEvent {
			whereText += " EventId >= ? "
		} else {
			whereText += " EventId > ? "
		}
		sqlParams = append(sqlParams, filterFromEventId)
	}

	if len(filterEntityType) > 0 {
		appendWhereOrAndToSqlText()
		whereText += "EntityType IN (?)"
		sqlParams = append(sqlParams, filterEntityType)
	}

	if len(filterEventType) > 0 {
		appendWhereOrAndToSqlText()
		whereText += "EventType IN (?)"
		sqlParams = append(sqlParams, filterEventType)
	}

	sqlText += whereText + " ORDER BY EventId ASC LIMIT " + rowsPerLoop

	rows, queryErr := cr.dbConnect.Query(sqlText, sqlParams...)
	if nil != queryErr {
		return nil, queryErr
	}

	return rows, nil
}

func NewClickhouseEventRepository(dbConnect *sql.DB, opentracingBridge *opentracingBridge.Bridge) EventRepository {
	return &clickhouseRepository{
		dbConnect:         dbConnect,
		opentracingBridge: opentracingBridge,
	}
}
