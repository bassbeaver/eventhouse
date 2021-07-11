package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bassbeaver/eventhouse/load-testing/api/compiled/event"
	"github.com/bojand/ghz/printer"
	"github.com/bojand/ghz/runner"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

var apiClientCredentials = []string{
	"Y2xpZW50MTpzZWNyZXQx", // client1:secret1
	"Y2xpZW50MjpzZWNyZXQy", // client2:secret2
}

var entityTypes = []string{"EntityType1", "EntityType2", "EntityType3", "EntityType4", "EntityType5"}
var createdEntities = NewSelfCleaningConcurrentSlice(500)

func main() {
	var target string

	flags := flag.NewFlagSet("flags", flag.PanicOnError)
	targetFlag := flags.String("target", "", "host and port of service to be tested")
	pathToProtoFlag := flags.String("proto", "", "path to .proto file describing api")
	flagsErr := flags.Parse(os.Args[1:])
	if nil != flagsErr {
		panic(flagsErr)
	}

	target = *targetFlag
	if "" == target {
		target = "localhost:750"
	}

	pathToProto := *pathToProtoFlag
	if "" == pathToProto {
		pathToProto = "/app/src/api/proto/event.proto"
	}

	fmt.Printf("Starting load test on %s \n", target)

	runs := [2]*RunAct{}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		fmt.Println(time.Now().Format(time.RFC3339) + " push started")
		runs[0] = runPushRunner(target, pathToProto)
		wg.Done()
	}()

	// To run stream reading test without parallel writing - comment out "readEntitiesFromGlobal(target)" line and "time.Sleep(1 * time.Minute)" line
	//readEntitiesFromGlobal(target)
	wg.Add(1)
	go func() {
		//time.Sleep(1 * time.Minute) // While testing parallel writing and reading - wait a while before reading in order to write routine to create some events
		fmt.Println(time.Now().Format(time.RFC3339) + " stream reading started")
		runs[1] = runEntityStreamRunner(target, pathToProto)
		wg.Done()
	}()

	wg.Wait()

	printReport(runs[:])
}

func authMetadataProvider(_ *runner.CallData) (*metadata.MD, error) {
	creds := apiClientCredentials[rand.Intn(len(apiClientCredentials))]
	md := metadata.New(map[string]string{"Authorization": "Basic " + creds})

	return &md, nil
}

func printReport(runs []*RunAct) {
	// Counting not empty reports
	var runsCount int
	for _, run := range runs {
		if nil != run {
			runsCount++
		}
	}
	fmt.Printf("Printing %d reports......\n\n", runsCount)

	for _, run := range runs {
		if nil == run {
			continue
		}

		if nil != run.Err {
			fmt.Printf("Run failed: %s \n\n", run.Err.Error())

			return
		}

		stars := strings.Repeat("*", len(run.Report.Name))
		fmt.Printf("%s\n %s \n%s", stars, run.Report.Name, stars)

		pushPrinterObj := printer.ReportPrinter{
			Out:    os.Stdout,
			Report: run.Report,
		}

		printErr := pushPrinterObj.Print("summary")
		if nil != printErr {
			fmt.Printf("Run result printing failed: %s \n", printErr.Error())
		}
	}
}

func readEntitiesFromGlobal(target string) {
	conn, err := grpc.Dial(
		target,
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(&ClientAuthCredentials{}),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		closeError := conn.Close()
		if nil != closeError {
			fmt.Printf("gRPC connection close error: %s", closeError.Error())
		}
	}()

	apiClient := event.NewAPIClient(conn)
	stream, streamError := apiClient.GlobalStream(context.Background(), &event.GlobalStreamRequest{})
	if nil != streamError {
		panic(streamError)
	}

	var readEvents int
	readEntities := make(map[string]bool)
	for {
		eventObj, err := stream.Recv()
		if err == io.EOF {
			break
		} else if nil != err {
			panic("error reading events stream: " + err.Error())
		}

		readEvents++

		entityKey := eventObj.GetEntityType() + eventObj.GetEntityId()
		if _, isRead := readEntities[entityKey]; isRead {
			continue
		}

		createdEntities.Append(
			&Entity{
				EntityType:    eventObj.GetEntityType(),
				EntityId:      eventObj.GetEntityId(),
				EventsCounter: NewMutexCounter(),
			},
		)

		readEntities[entityKey] = true
	}

	fmt.Printf("Read %d events. Created %d entities \n", readEvents, createdEntities.Len())
}
