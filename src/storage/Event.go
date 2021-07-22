package storage

import (
	"crypto/rand"
	"errors"
	"github.com/bassbeaver/logopher"
	"math/big"
	"time"
)

const (
	EventRepositoryAlias = "EventRepository"
	rowsPerLoop          = "10"
	streamChanBufferSize = 15 // 1.5 * rowsPerLoop
)

type Event struct {
	EventId    uint64
	EventType  string
	EntityType string
	EntityId   string
	Recorded   time.Time
	Payload    string
}

type EventRepository interface {
	Save(eventType string, idempotencyKey string, entityType string, entityId string, payload string) (*Event, error)
	Get(eventId uint64) (*Event, error)
	EntityStream(entityType, entityId string, filterFromEventId uint64, loggerObj *logopher.Logger) (chan *Event, error)
	GlobalStream(filterFromEventId uint64, filterEntityType, filterEventType []string, loggerObj *logopher.Logger) (chan *Event, error)
}

func generateNewEventId() (uint64, error) {
	randPart, err := rand.Int(rand.Reader, big.NewInt(9))
	if nil != err {
		return 0, errors.New("failed to generate random part of new uuid")
	}

	return uint64(time.Now().UnixNano())*10 + uint64(randPart.Int64()), nil
}
