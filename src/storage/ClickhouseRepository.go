package storage

import (
	"database/sql"
	"errors"
	"github.com/bassbeaver/logopher"
	"time"
)

type clickhouseRepository struct {
	dbConnect *sql.DB
}

func (cr *clickhouseRepository) Save(
	eventType string,
	idempotencyKey string,
	entityType string,
	entityId string,
	payload string,
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

	tx, txError := cr.dbConnect.Begin()
	if nil != txError {
		return nil, errors.New("Failed to create Clickhouse transaction: " + txError.Error())
	}

	stmt, prepareErr := tx.Prepare(
		`INSERT INTO events (EventId, EventType, IdempotencyKey, EntityType, EntityId, Recorded, Payload)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
	)
	if nil != prepareErr {
		return nil, errors.New("Failed to prepare insert statement: " + prepareErr.Error())
	}

	_, stmtExecErr := stmt.Exec(
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
		"SELECT EventId, EventType, EntityType, EntityId, Recorded, Payload FROM events WHERE EventId = ? LIMIT 1",
		eventId,
	)
	if nil != queryErr {
		return nil, errors.New("Failed to query Event: " + queryErr.Error())
	}

	if !rows.Next() {
		return nil, nil
	}

	event := &Event{}
	scanErr := rows.Scan(&event.EventId, &event.EventType, &event.EntityType, &event.EntityId, &event.Recorded, &event.Payload)
	if nil != scanErr {
		return nil, errors.New("Failed to scan query results: " + scanErr.Error())
	}

	return event, nil
}

func (cr *clickhouseRepository) EntityStream(
	entityType string,
	entityId string,
	filterFromEventId uint64,
	loggerObj *logopher.Logger,
) (chan *Event, error) {
	// Buffered channel used to avoid case when stream reader is very slow and repository loaded all events from DB and pushed it to chan.
	// We want to load events with same pace as they are read.
	eventsChan := make(chan *Event, streamChanBufferSize)

	go cr.performStreamRead(
		func(lastEventId uint64) (*sql.Rows, error) {
			return cr.performEntityStreamBatchQuery(entityType, entityId, filterFromEventId, lastEventId)
		},
		eventsChan,
		loggerObj,
	)

	return eventsChan, nil
}

func (cr *clickhouseRepository) GlobalStream(
	filterFromEventId uint64,
	filterEntityType []string,
	filterEventType []string,
	loggerObj *logopher.Logger,
) (chan *Event, error) {
	eventsChan := make(chan *Event, streamChanBufferSize)

	go cr.performStreamRead(
		func(lastEventId uint64) (*sql.Rows, error) {
			return cr.performGlobalStreamBatchQuery(filterFromEventId, filterEntityType, filterEventType, lastEventId)
		},
		eventsChan,
		loggerObj,
	)

	return eventsChan, nil
}

func (cr *clickhouseRepository) performStreamRead(
	batchQueryPerformer func(filterFromEventId uint64) (*sql.Rows, error),
	eventsChan chan *Event,
	loggerObj *logopher.Logger,
) {
	var lastEventId uint64
	var rows *sql.Rows
	queryNextBatch := true

	for queryNextBatch {
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
			queryNextBatch = true

			event := &Event{}
			scanErr := rows.Scan(&event.EventId, &event.EventType, &event.EntityType, &event.EntityId, &event.Recorded, &event.Payload)
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
	lastEventId uint64,
) (*sql.Rows, error) {
	sqlText :=
		`SELECT EventId, EventType, EntityType, EntityId, Recorded, Payload FROM events
		WHERE EntityType = ? AND EntityId = ? `
	sqlParams := []interface{}{entityType, entityId}

	if 0 != lastEventId {
		sqlText += " AND EventId > ? "
		sqlParams = append(sqlParams, lastEventId)
	}

	if 0 != filterFromEventId {
		sqlText += " AND EventId >= ? "
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
	filterEntityType []string,
	filterEventType []string,
	lastEventId uint64,
) (*sql.Rows, error) {
	sqlText := "SELECT EventId, EventType, EntityType, EntityId, Recorded, Payload FROM events"
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
		whereText += " EventId >= ? "
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

func NewClickhouseEventRepository(dbConnect *sql.DB) EventRepository {
	return &clickhouseRepository{dbConnect: dbConnect}
}
