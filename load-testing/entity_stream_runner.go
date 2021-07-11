package main

import (
	"fmt"
	"github.com/bassbeaver/eventhouse/load-testing/api/compiled/event"
	"github.com/bojand/ghz/runner"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"math/rand"
)

func runEntityStreamRunner(target, proto string) *RunAct {
	act := &RunAct{}

	act.Report, act.Err = runner.Run(
		"eventhouse.grpc.event.API.EntityStream",
		target,
		runner.WithName("eventhouse.grpc.event.API.EntityStream"),
		runner.WithProtoFile(proto, []string{}),
		runner.WithBinaryDataFunc(getDataProvider),
		runner.WithMetadataProvider(authMetadataProvider),
		runner.WithInsecure(true),
		runner.WithConcurrency(10),
		runner.WithTotalRequests(15000),
		runner.WithLoadSchedule(runner.ScheduleLine),
		runner.WithLoadStart(2),
		runner.WithLoadEnd(50),
		runner.WithLoadStep(2),
	)

	return act
}

func getDataProvider(_ *desc.MethodDescriptor, _ *runner.CallData) []byte {
	if createdEntities.Len() == 0 {
		panic("Failed to start stream reading test. No entities created.")
	}

	i := rand.Intn(createdEntities.Len())
	entityIntf, _ := createdEntities.Get(i)
	entity := entityIntf.(*Entity)

	msg := &event.EntityStreamRequest{}
	msg.EntityId = entity.EntityId
	msg.EntityType = entity.EntityType

	binData, err := proto.Marshal(msg)
	if nil != err {
		fmt.Printf("Error in pushDataProvider: %s \n", err)
	}

	return binData
}
