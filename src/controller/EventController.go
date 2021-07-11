package controller

import (
	"context"
	apiEvent "github.com/bassbeaver/eventhouse/api/compiled/event"
	loggerService "github.com/bassbeaver/eventhouse/service/logger"
	"github.com/bassbeaver/eventhouse/storage"
	"github.com/bassbeaver/logopher"
	grpcCodes "google.golang.org/grpc/codes"
	grpcStatus "google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"strconv"
)

const (
	EventControllerServiceAlias = "EventController"
)

type EventController struct {
	eventRepo storage.EventRepository
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

	eventsChan, chanError := ec.eventRepo.EntityStream(requestObj.GetEntityType(), requestObj.GetEntityId(), fromEventId, nil)
	if nil != chanError {
		loggerObj.Critical("Failed create Events stream channel", &logopher.MessageContext{"error": chanError.Error()})

		return grpcStatus.Error(grpcCodes.Internal, "server error")
	}

	for evt := range eventsChan {
		sendError := streamServer.Send(dbToApi(evt))
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

	eventsChan, chanError := ec.eventRepo.GlobalStream(fromEventId, requestObj.GetEntityType(), requestObj.GetEventType(), nil)
	if nil != chanError {
		loggerObj.Critical("Failed create Events stream channel", &logopher.MessageContext{"error": chanError.Error()})

		return grpcStatus.Error(grpcCodes.Internal, "server error")
	}

	for evt := range eventsChan {
		sendError := streamServer.Send(dbToApi(evt))
		if nil != sendError {
			loggerObj.Critical("Failed send Event to stream", &logopher.MessageContext{"error": sendError.Error()})

			return grpcStatus.Error(grpcCodes.Internal, "server error")
		}
	}

	return nil
}

// --------

func NewEventController(eventRepo storage.EventRepository) *EventController {
	return &EventController{eventRepo: eventRepo}
}

func dbToApi(event *storage.Event) *apiEvent.Event {
	return &apiEvent.Event{
		EventId:    strconv.FormatUint(event.EventId, 10),
		EventType:  event.EventType,
		EntityType: event.EntityType,
		EntityId:   event.EntityId,
		Recorded:   timestamppb.New(event.Recorded),
		Payload:    event.Payload,
	}
}
