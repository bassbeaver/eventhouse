package storage

import (
	"context"
	"database/sql"
	"fmt"
	opentracingBridge "github.com/bassbeaver/eventhouse/service/opentracing"
	"github.com/opentracing/opentracing-go"
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
		Payload:    payload,
	}

	saveResultChannel := make(chan error)

	if nil != ctx {
		if parentSpan := opentracing.SpanFromContext(ctx); nil != parentSpan {
			repoSaveSpan := bcr.opentracingBridge.Tracer().StartSpan(
				"event_repo__save",
				opentracing.ChildOf(parentSpan.Context()),
				opentracing.Tag{Key: "EventId", Value: fmt.Sprintf("%d", newEvent.EventId)},
			)
			defer repoSaveSpan.Finish()
		}
	}

	bcr.batchManagerObj.Append(newBatchEntity(newEvent, idempotencyKey, saveResultChannel))

	if saveError := <-saveResultChannel; nil != saveError {
		return nil, saveError
	}

	return newEvent, nil
}

func NewBatchSavingClickhouseRepository(
	maxEntitiesInBatch int,
	batchLifetimeMs int,
	dbConnect *sql.DB,
	opentracingBridge *opentracingBridge.Bridge,
) EventRepository {
	return &batchSavingClickhouseRepository{
		clickhouseRepository: NewClickhouseEventRepository(dbConnect, opentracingBridge).(*clickhouseRepository),
		batchManagerObj:      newBatchManager(maxEntitiesInBatch, time.Duration(batchLifetimeMs)*time.Millisecond, dbConnect),
	}
}
