package storage

import (
	"database/sql"
	"time"
)

type batchSavingClickhouseRepository struct {
	*clickhouseRepository
	batchManagerObj *batchManager
}

func (bcr *batchSavingClickhouseRepository) Save(
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

	saveResultChannel := make(chan error)

	bcr.batchManagerObj.Append(newBatchEntity(newEvent, idempotencyKey, saveResultChannel))

	if saveError := <-saveResultChannel; nil != saveError {
		return nil, saveError
	}

	return newEvent, nil
}

func NewBatchSavingClickhouseRepository(maxEntitiesInBatch int, batchLifetimeMs int, dbConnect *sql.DB) EventRepository {
	return &batchSavingClickhouseRepository{
		clickhouseRepository: NewClickhouseEventRepository(dbConnect).(*clickhouseRepository),
		batchManagerObj:      newBatchManager(maxEntitiesInBatch, time.Duration(batchLifetimeMs)*time.Millisecond, dbConnect),
	}
}
