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
	entitiesCount      int
	maxEntities        int
	maxLifeTime        time.Duration
	insertParams       [][7]interface{}
	saveResultChannels []chan error
	dbConnect          *sql.DB
	appendChan         chan *batchEntity
	closedChan         chan bool
	lifetimeTicker     *time.Ticker
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
				b.saveResultChannels = append(b.saveResultChannels, batchEntityObj.saveResultChannel)
				b.insertParams = append(
					b.insertParams,
					[7]interface{}{
						batchEntityObj.event.EventId,
						batchEntityObj.event.EventType,
						batchEntityObj.idempotencyKey,
						batchEntityObj.event.EntityType,
						batchEntityObj.event.EntityId,
						batchEntityObj.event.Recorded,
						batchEntityObj.event.Payload,
					},
				)

				if b.entitiesCount >= b.maxEntities {
					b.performInsertAndSendResults()

					return // stopping current batch
				}
			case <-b.lifetimeTicker.C:
				// If no values were submitted for save - just wait for them
				if len(b.insertParams) == 0 {
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
		for _, saveResultChannel := range b.saveResultChannels {
			saveResultChannel <- insertError
		}
	} else {
		for _, saveResultChannel := range b.saveResultChannels {
			saveResultChannel <- nil
		}
	}
}

func (b *batch) performInsert() error {
	tx, txError := b.dbConnect.Begin()
	if nil != txError {
		return errors.New("Failed to create Clickhouse transaction: " + txError.Error())
	}

	stmt, prepareErr := tx.Prepare(
		"INSERT INTO events (EventId, EventType, IdempotencyKey, EntityType, EntityId, Recorded, Payload) VALUES (?, ?, ?, ?, ?, ?, ?)",
	)
	if nil != prepareErr {
		return errors.New("Failed to prepare insert statement: " + prepareErr.Error())
	}
	for _, rowParams := range b.insertParams {
		_, stmtExecErr := stmt.Exec(rowParams[:]...)
		if nil != stmtExecErr {
			return errors.New("Failed to execute prepared insert statement: " + stmtExecErr.Error())
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
		maxEntities:        maxEntities,
		maxLifeTime:        maxLifeTime,
		insertParams:       make([][7]interface{}, 0),
		saveResultChannels: make([]chan error, 0),
		dbConnect:          dbConnect,
		appendChan:         appendChan,
		closedChan:         make(chan bool),
	}

	b.startAppendListener()

	return b
}
