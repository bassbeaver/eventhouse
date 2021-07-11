package main

import (
	"fmt"
	"github.com/bassbeaver/eventhouse/load-testing/api/compiled/event"
	"github.com/bojand/ghz/runner"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"math/rand"
	"strconv"
	"time"
)

const maxEventsForEntity = 15

func runPushRunner(target, proto string) *RunAct {
	act := &RunAct{}

	act.Report, act.Err = runner.Run(
		"eventhouse.grpc.event.API.Push",
		target,
		runner.WithName("eventhouse.grpc.event.API.Push"),
		runner.WithProtoFile(proto, []string{}),
		runner.WithBinaryDataFunc(pushDataProvider),
		runner.WithMetadataProvider(authMetadataProvider),
		runner.WithInsecure(true),
		runner.WithConcurrency(10),
		runner.WithTotalRequests(10000),
		runner.WithLoadSchedule(runner.ScheduleStep),
		runner.WithLoadStart(10),
		runner.WithLoadEnd(40),
		runner.WithLoadStep(5),
		runner.WithLoadStepDuration(10*time.Second),
	)

	return act
}

func pushDataProvider(_ *desc.MethodDescriptor, _ *runner.CallData) []byte {
	var entity *Entity

	if createdEntities.Len() > 10 {
		i := rand.Intn(createdEntities.Len())
		entityIntf, _ := createdEntities.Get(i)
		entity = entityIntf.(*Entity)

		// If Entity already had maxEventsForEntity events - create new
		eventsCount := entity.EventsCounter.Current()
		if eventsCount >= maxEventsForEntity {
			entity = nil
		}
	}

	if nil == entity {
		entity = newEntity()
		createdEntities.Append(entity)
	}

	idempotencyKeyBytes := make([]byte, 16)
	rand.Read(idempotencyKeyBytes)
	eventsCount := entity.EventsCounter.Increment()
	eventType := fmt.Sprintf("event-%d", eventsCount)
	msg := &event.PushRequest{
		IdempotencyKey: fmt.Sprintf("%x", idempotencyKeyBytes),
		EntityType:     entity.EntityType,
		EntityId:       entity.EntityId,
		EventType:      eventType,
		Payload:        fmt.Sprintf(`{"data": "payload-for-%s", "entity": "%s--%s"}`, eventType, entity.EntityType, entity.EntityId),
	}

	binData, err := proto.Marshal(msg)
	if nil != err {
		fmt.Printf("Error in pushDataProvider: %s \n", err)
	}

	return binData
}

func newEntity() *Entity {
	return &Entity{
		EntityType:    entityTypes[rand.Intn(len(entityTypes))],
		EntityId:      "id-" + strconv.Itoa(createdEntities.Total()),
		EventsCounter: NewMutexCounter(),
	}
}
