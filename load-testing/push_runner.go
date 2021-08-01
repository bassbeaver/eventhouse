package main

import (
	"fmt"
	"github.com/bojand/ghz/runner"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"math/rand"
	"time"
)

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
		runner.WithTotalRequests(2000),
		runner.WithLoadSchedule(runner.ScheduleStep),
		runner.WithLoadStart(10),
		runner.WithLoadEnd(50),
		runner.WithLoadStep(2),
		runner.WithLoadStepDuration(10*time.Second),
	)

	return act
}

func pushDataProvider(_ *desc.MethodDescriptor, _ *runner.CallData) []byte {
	seq := getActiveSubscriptionSequenceAndGenerateNextEvent(0)
	evtObj := seq.GetLastEvent()
	binData, err := proto.Marshal(evtObj)
	if nil != err {
		fmt.Printf("Error in pushDataProvider: %s \n", err)
	}

	loseDice := rand.Intn(100)
	if loseDice < 85 { // 15% chance to lose subscription and left it in active state
		subscriptionRegistry.FreeSequence(seq)
	}

	return binData
}

func getActiveSubscriptionSequenceAndGenerateNextEvent(loopDepth int) *SubscriptionSequence {
	seq := subscriptionRegistry.GetSequence()

	if !seq.IsCanceled() && seq.GenerateNextEvent() {
		return seq
	}

	if loopDepth >= 10 {
		panic(fmt.Sprintf("failed to found active SubscriptionSequence during %d loops", loopDepth))
	}

	return getActiveSubscriptionSequenceAndGenerateNextEvent(loopDepth + 1)
}
