package storage

import (
	"database/sql"
	"errors"
	"time"
)

type batchEntity struct {
	event             *Event
	idempotencyKey    string
	saveResultChannel chan error
}

func newBatchEntity(event *Event, idempotencyKey string, saveResultChannel chan error) *batchEntity {
	return &batchEntity{
		event:             event,
		idempotencyKey:    idempotencyKey,
		saveResultChannel: saveResultChannel,
	}
}

type batch struct {
	entitiesCount         int
	maxEntities           int
	maxLifeTime           time.Duration
	batchEntitiesToInsert []*batchEntity
	//saveResultChannels    []chan error
	dbConnect      *sql.DB
	appendChan     chan *batchEntity
	closedChan     chan bool
	lifetimeTicker *time.Ticker
}

func (b *batch) startAppendListener() {
	b.lifetimeTicker = time.NewTicker(b.maxLifeTime)

	go func() {
		for {
			select {
			case batchEntityObj, writeChanIsOk := <-b.appendChan:
				if !writeChanIsOk {
					return
				}

				b.entitiesCount++
				//b.saveResultChannels = append(b.saveResultChannels, batchEntityObj.saveResultChannel)
				b.batchEntitiesToInsert = append(b.batchEntitiesToInsert, batchEntityObj)

				if b.entitiesCount >= b.maxEntities {
					b.performInsertAndSendResults()

					return // stopping current batch
				}
			case <-b.lifetimeTicker.C:
				// If no values were submitted for save - just wait for them
				if len(b.batchEntitiesToInsert) == 0 {
					continue
				}

				// If values were submitted and lifetime spent - save values and stop current batch
				b.performInsertAndSendResults()

				return
			}
		}
	}()
}

func (b *batch) performInsertAndSendResults() {
	b.lifetimeTicker.Stop()
	b.closedChan <- true

	insertError := b.performInsert()
	if nil != insertError {
		//for _, saveResultChannel := range b.saveResultChannels {
		for _, entityObj := range b.batchEntitiesToInsert {
			entityObj.saveResultChannel <- insertError
		}
	} else {
		//for _, saveResultChannel := range b.saveResultChannels {
		for _, entityObj := range b.batchEntitiesToInsert {
			entityObj.saveResultChannel <- nil
		}
	}
}

func (b *batch) performInsert() error {
	tx, txError := b.dbConnect.Begin()
	if nil != txError {
		return errors.New("Failed to create Clickhouse transaction: " + txError.Error())
	}

	// Checking idempotency
	alreadyExistingIdempotencyKeys := make(map[string]bool)
	idempotencyKeys := make([]string, len(b.batchEntitiesToInsert))
	for i, batchEntityObj := range b.batchEntitiesToInsert {
		idempotencyKeys[i] = batchEntityObj.idempotencyKey
	}
	alreadyExistingIdempotencyKeysRows, alreadyExistingQueryError := tx.Query(
		"SELECT IdempotencyKey FROM events WHERE IdempotencyKey IN (?)",
		idempotencyKeys,
	)
	if nil != alreadyExistingQueryError {
		return errors.New("Failed to execute idempotency check query: " + alreadyExistingQueryError.Error())
	}
	defer alreadyExistingIdempotencyKeysRows.Close()
	for alreadyExistingIdempotencyKeysRows.Next() {
		var k string
		if scanError := alreadyExistingIdempotencyKeysRows.Scan(&k); scanError != nil {
			return errors.New("Failed to scan idempotency check query result: " + scanError.Error())
		}

		alreadyExistingIdempotencyKeys[k] = true
	}

	// Inserting events
	// NOTICE:
	// In select statement subquery have to be at first place, because database/sql package has a bug and
	// do not recognize prepared statement placeholder in SELECT statement if it is located just after SELECT keyword
	stmt, prepareErr := tx.Prepare(
		"INSERT INTO events (PreviousEventId, EventId, EventType, IdempotencyKey, EntityType, EntityId, Recorded, Payload) " +
			"SELECT " +
			"(SELECT IF(0 < MAX(EventId), MAX(EventId), 0) FROM events WHERE EntityType=? AND EntityId=?) as PreviousEventId, " +
			"? as EventId, " +
			"? as EventType, " +
			"? as IdempotencyKey, " +
			"? as EntityType, " +
			"? as EntityId, " +
			"? as Recorded, " +
			"? as Payload",
	)
	if nil != prepareErr {
		return errors.New("Failed to prepare insert statement: " + prepareErr.Error())
	}
	for _, batchEntityObj := range b.batchEntitiesToInsert {
		if alreadyExistingIdempotencyKeys[batchEntityObj.idempotencyKey] {
			continue
		}

		_, stmtExecErr := stmt.Exec(
			batchEntityObj.event.EntityType,
			batchEntityObj.event.EntityId,
			batchEntityObj.event.EventId,
			batchEntityObj.event.EventType,
			batchEntityObj.idempotencyKey,
			batchEntityObj.event.EntityType,
			batchEntityObj.event.EntityId,
			batchEntityObj.event.Recorded,
			batchEntityObj.event.Payload,
		)
		if nil != stmtExecErr {
			return errors.New("Failed to execute prepared insert statement: " + stmtExecErr.Error())
		}
	}

	// Querying PreviousEventId for events
	previousIds := make(map[uint64]uint64) // map[EventId]PreviousEventId
	eventIds := make([]uint64, len(b.batchEntitiesToInsert))
	for i, batchEntityObj := range b.batchEntitiesToInsert {
		eventIds[i] = batchEntityObj.event.EventId
	}
	eventsRows, eventsQueryError := tx.Query(
		"SELECT EventId, PreviousEventId FROM events WHERE EventId IN (?)",
		eventIds,
	)
	if nil != eventsQueryError {
		return errors.New("Failed to execute previous events ids query: " + eventsQueryError.Error())
	}
	defer eventsRows.Close()
	for eventsRows.Next() {
		var id, prevId uint64
		if scanError := eventsRows.Scan(&id, &prevId); scanError != nil {
			return errors.New("Failed to scan previous events ids query result: " + scanError.Error())
		}

		previousIds[id] = prevId
	}

	// Setting PreviousEventId to events
	for _, batchEntityObj := range b.batchEntitiesToInsert {
		if previousId, previousIdIsOk := previousIds[batchEntityObj.event.EventId]; previousIdIsOk {
			batchEntityObj.event.PreviousEventId = previousId
		}
	}

	commitErr := tx.Commit()
	if nil != commitErr {
		return errors.New("Failed to commit transaction: " + commitErr.Error())
	}

	return nil
}

func newBatch(
	maxEntities int,
	maxLifeTime time.Duration,
	dbConnect *sql.DB,
	appendChan chan *batchEntity,
) *batch {
	b := &batch{
		maxEntities:           maxEntities,
		maxLifeTime:           maxLifeTime,
		batchEntitiesToInsert: make([]*batchEntity, 0),
		//saveResultChannels:    make([]chan error, 0),
		dbConnect:  dbConnect,
		appendChan: appendChan,
		closedChan: make(chan bool),
	}

	b.startAppendListener()

	return b
}
