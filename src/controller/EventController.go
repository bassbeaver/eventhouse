package controller

import (
	"context"
	apiEvent "github.com/bassbeaver/eventhouse/api/compiled/event"
	loggerService "github.com/bassbeaver/eventhouse/service/logger"
	opentracingBridge "github.com/bassbeaver/eventhouse/service/opentracing"
	"github.com/bassbeaver/eventhouse/storage"
	"github.com/bassbeaver/logopher"
	"github.com/opentracing/opentracing-go"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
	"time"
)

const (
	EventControllerServiceAlias = "EventController"
)

type EventController struct {
	eventRepo         storage.EventRepository
	opentracingBridge *opentracingBridge.Bridge
}

func (ec *EventController) Push(ctx context.Context, requestObj *apiEvent.PushRequest) (*apiEvent.Event, error) {
	loggerObj := loggerService.GetLoggerFromContext(ctx)

	if "" == requestObj.EventType {
		loggerObj.Warning("Empty Event Type not allowed", nil)

		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "invalid argument")
	}

	newEvent, eventSaveError := ec.eventRepo.Save(
		requestObj.GetEventType(),
		requestObj.GetIdempotencyKey(),
		requestObj.GetEntityType(),
		requestObj.GetEntityId(),
		requestObj.GetPayload(),
		ctx,
	)

	if nil != eventSaveError {
		loggerObj.Critical("Failed to save Event to Repo", &logopher.MessageContext{"error": eventSaveError.Error()})

		return nil, grpcStatus.Error(grpcCodes.Internal, "server error")
	}

	return dbToApi(newEvent), nil
}

func (ec *EventController) Get(ctx context.Context, requestObj *apiEvent.GetRequest) (*apiEvent.Event, error) {
	loggerObj := loggerService.GetLoggerFromContext(ctx)

	eventId, eventIdParseError := strconv.ParseUint(requestObj.GetEventId(), 10, 64)
	if nil != eventIdParseError {
		loggerObj.Warning("Failed to parse eventId to uint64", &logopher.MessageContext{"error": eventIdParseError.Error()})

		return nil, grpcStatus.Error(grpcCodes.InvalidArgument, "invalid argument")
	}

	event, eventGetError := ec.eventRepo.Get(eventId)
	if nil != eventGetError {
		loggerObj.Critical("Failed to get Event from repo", &logopher.MessageContext{"error": eventGetError.Error()})

		return nil, grpcStatus.Error(grpcCodes.Internal, "server error")
	}

	if nil == event {
		return nil, grpcStatus.Error(grpcCodes.NotFound, "event not found")
	}

	return dbToApi(event), nil
}

func (ec *EventController) EntityStream(requestObj *apiEvent.EntityStreamRequest, streamServer apiEvent.API_EntityStreamServer) error {
	var fromEventId uint64
	var eventIdParseError error

	loggerObj := loggerService.GetLoggerFromContext(streamServer.Context())

	if "" != requestObj.GetFilter().GetEventIdFrom() {
		fromEventId, eventIdParseError = strconv.ParseUint(requestObj.GetFilter().GetEventIdFrom(), 10, 64)
		if nil != eventIdParseError {
			loggerObj.Warning("Failed to parse EventId to uint64", &logopher.MessageContext{"error": eventIdParseError.Error()})

			return grpcStatus.Error(grpcCodes.InvalidArgument, "invalid argument")
		}
	}

	opentracingRootSpan := opentracing.SpanFromContext(streamServer.Context())

	eventsChan, chanError := ec.eventRepo.EntityStream(
		requestObj.GetEntityType(),
		requestObj.GetEntityId(),
		fromEventId,
		true,
		nil,
		opentracing.ContextWithSpan(context.Background(), opentracingRootSpan),
	)
	if nil != chanError {
		loggerObj.Critical("Failed create Events stream channel", &logopher.MessageContext{"error": chanError.Error()})

		return grpcStatus.Error(grpcCodes.Internal, "server error")
	}

	for evt := range eventsChan {
		childSpan := ec.opentracingBridge.Tracer().StartSpan(
			"event_send",
			opentracing.ChildOf(opentracingRootSpan.Context()),
		)

		quantumMeta := make(opentracing.TextMapCarrier)
		injectError := ec.opentracingBridge.Tracer().Inject(childSpan.Context(), opentracing.TextMap, quantumMeta)
		if injectError != nil {
			loggerObj.Error("Failed to inject Opentracing data to EventStreamQuantum meta", &logopher.MessageContext{"error": injectError.Error()})
		}

		sendError := streamServer.Send(&apiEvent.EventStreamQuantum{Event: dbToApi(evt), Meta: quantumMeta})

		childSpan.Finish()

		if nil != sendError {
			loggerObj.Critical("Failed send Event to stream", &logopher.MessageContext{"error": sendError.Error()})

			return grpcStatus.Error(grpcCodes.Internal, "server error")
		}
	}

	return nil
}

func (ec *EventController) GlobalStream(requestObj *apiEvent.GlobalStreamRequest, streamServer apiEvent.API_GlobalStreamServer) error {
	var fromEventId uint64
	var eventIdParseError error

	loggerObj := loggerService.GetLoggerFromContext(streamServer.Context())

	if "" != requestObj.GetEventIdFrom() {
		fromEventId, eventIdParseError = strconv.ParseUint(requestObj.GetEventIdFrom(), 10, 64)
		if nil != eventIdParseError {
			loggerObj.Warning("Failed to parse EventId to uint64", &logopher.MessageContext{"error": eventIdParseError.Error()})

			return grpcStatus.Error(grpcCodes.InvalidArgument, "invalid argument")
		}
	}

	opentracingRootSpan := opentracing.SpanFromContext(streamServer.Context())

	eventsChan, chanError := ec.eventRepo.GlobalStream(
		fromEventId,
		true,
		requestObj.GetEntityType(),
		requestObj.GetEventType(),
		loggerObj,
		opentracing.ContextWithSpan(context.Background(), opentracingRootSpan),
	)
	if nil != chanError {
		loggerObj.Critical("Failed create Events stream channel", &logopher.MessageContext{"error": chanError.Error()})

		return grpcStatus.Error(grpcCodes.Internal, "server error")
	}

	for evt := range eventsChan {
		childSpan := ec.opentracingBridge.Tracer().StartSpan(
			"controller__event_send",
			opentracing.ChildOf(opentracingRootSpan.Context()),
			opentracing.Tag{Key: "EventId", Value: strconv.FormatUint(evt.EventId, 10)},
		)

		quantumMeta := make(opentracing.TextMapCarrier)
		injectError := ec.opentracingBridge.Tracer().Inject(childSpan.Context(), opentracing.TextMap, quantumMeta)
		if injectError != nil {
			loggerObj.Error("Failed to inject Opentracing data to EventStreamQuantum meta", &logopher.MessageContext{"error": injectError.Error()})
		}

		sendError := streamServer.Send(&apiEvent.EventStreamQuantum{Event: dbToApi(evt), Meta: quantumMeta})

		childSpan.Finish()

		if nil != sendError {
			loggerObj.Critical("Failed send Event to stream", &logopher.MessageContext{"error": sendError.Error()})

			return grpcStatus.Error(grpcCodes.Internal, "server error")
		}
	}

	return nil
}

func (ec *EventController) SubscribeGlobalStream(requestObj *apiEvent.SubscribeGlobalStreamRequest, streamServer apiEvent.API_SubscribeGlobalStreamServer) error {
	var fromEventId uint64
	var eventIdParseError error

	loggerObj := loggerService.GetLoggerFromContext(streamServer.Context())

	if "" != requestObj.GetEventIdFrom() {
		fromEventId, eventIdParseError = strconv.ParseUint(requestObj.GetEventIdFrom(), 10, 64)
		if nil != eventIdParseError {
			loggerObj.Warning("Failed to parse EventId to uint64", &logopher.MessageContext{"error": eventIdParseError.Error()})

			return grpcStatus.Error(grpcCodes.InvalidArgument, "invalid argument")
		}
	}

	// Read already saved events

	eventsChan, chanError := ec.eventRepo.GlobalStream(
		fromEventId,
		true,
		requestObj.GetEntityType(),
		requestObj.GetEventType(),
		loggerObj,
		streamServer.Context(),
	)
	if nil != chanError {
		loggerObj.Critical("Failed create Events stream channel", &logopher.MessageContext{"error": chanError.Error()})

		return grpcStatus.Error(grpcCodes.Internal, "server error")
	}

	for evt := range eventsChan {
		sendError := ec.sendEventToSubscription(evt, streamServer, loggerObj)
		if nil != sendError {
			loggerObj.Critical("Failed send Event to stream", &logopher.MessageContext{"error": sendError.Error()})

			return grpcStatus.Error(grpcCodes.Internal, "server error")
		}

		fromEventId = evt.EventId
	}

	// Start listening to new events

	ticker := time.NewTicker(1000 * time.Millisecond)
	done := make(chan error)
	go func() {
		for {
			select {
			case <-ticker.C:
				newEventsChan, newEventsChanError := ec.eventRepo.GlobalStream(
					fromEventId,
					false,
					requestObj.GetEntityType(),
					requestObj.GetEventType(),
					loggerObj,
					streamServer.Context(),
				)
				if nil != newEventsChanError {
					loggerObj.Critical("Failed create new Events stream channel", &logopher.MessageContext{"error": newEventsChanError.Error()})

					done <- grpcStatus.Error(grpcCodes.Internal, "server error")

					return
				}

				for evt := range newEventsChan {
					sendError := ec.sendEventToSubscription(evt, streamServer, loggerObj)
					if nil != sendError {
						loggerObj.Critical("Failed send new Event to stream", &logopher.MessageContext{"error": sendError.Error()})

						done <- grpcStatus.Error(grpcCodes.Internal, "server error")

						return
					}

					fromEventId = evt.EventId
				}
			case <-streamServer.Context().Done():
				done <- nil
			}
		}
	}()

	listenError := <-done
	if nil == listenError {
		loggerObj.Info("Events subscription stopped. StreamServer context is done.", nil)
	} else {
		loggerObj.Error("Events subscription stopped due to error.", &logopher.MessageContext{"error": listenError.Error()})
	}

	return listenError
}

func (ec *EventController) sendEventToSubscription(
	evt *storage.Event,
	streamServer apiEvent.API_SubscribeGlobalStreamServer,
	loggerObj *logopher.Logger,
) error {
	opentracingRootSpan := opentracing.SpanFromContext(streamServer.Context())

	rootSpan := ec.opentracingBridge.Tracer().StartSpan(
		"controller__subscription__event_send",
		opentracing.ChildOf(opentracingRootSpan.Context()),
		opentracing.Tag{Key: "EventId", Value: strconv.FormatUint(evt.EventId, 10)},
	)
	defer rootSpan.Finish()

	quantumMeta := make(opentracing.TextMapCarrier)
	injectError := ec.opentracingBridge.Tracer().Inject(rootSpan.Context(), opentracing.TextMap, quantumMeta)
	if injectError != nil {
		loggerObj.Error("Failed to inject Opentracing data to EventStreamQuantum meta", &logopher.MessageContext{"error": injectError.Error()})
	}

	return streamServer.Send(&apiEvent.EventStreamQuantum{Event: dbToApi(evt), Meta: quantumMeta})
}

// --------

func NewEventController(eventRepo storage.EventRepository, opentracingBridge *opentracingBridge.Bridge) *EventController {
	return &EventController{
		eventRepo:         eventRepo,
		opentracingBridge: opentracingBridge,
	}
}

func dbToApi(event *storage.Event) *apiEvent.Event {
	var previousEventId string
	if event.PreviousEventId > 0 {
		previousEventId = strconv.FormatUint(event.PreviousEventId, 10)
	}

	return &apiEvent.Event{
		EventId:         strconv.FormatUint(event.EventId, 10),
		EventType:       event.EventType,
		EntityType:      event.EntityType,
		EntityId:        event.EntityId,
		Recorded:        timestamppb.New(event.Recorded),
		Payload:         event.Payload,
		PreviousEventId: previousEventId,
	}
}
