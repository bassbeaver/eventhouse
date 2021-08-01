package event_reader_reducer

import (
	"context"
	"fmt"
	apiEvent "github.com/bassbeaver/eventhouse/api/compiled/event"
	projectorOpentracing "github.com/bassbeaver/eventhouse/projector/service/opentracing"
	"github.com/opentracing/opentracing-go"
	"io"
)

const (
	EventReaderReducerServiceAlias = "EventReaderReducerService"
)

type EventReaderReducer struct {
	client  apiEvent.APIClient
	tracer  opentracing.Tracer
	reducer *EventReducer
}

func (er *EventReaderReducer) ReadReducerGlobalStream() {
	streamObj, streamError := er.client.GlobalStream(context.Background(), &apiEvent.GlobalStreamRequest{})
	if nil != streamError {
		panic("failed to call GlobalStream: " + streamError.Error())
	}

	streamClientSpan := projectorOpentracing.SpanFromContext(streamObj.Context())

	for {
		eventStreamQuantum, err := streamObj.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			panic("error while reading global stream: " + err.Error())
		}

		var eventSpan opentracing.Span
		if parentSpanCtx := er.selectOpentracingContextForEventStreamQuantum(streamClientSpan, eventStreamQuantum); nil != parentSpanCtx {
			eventSpan = er.tracer.StartSpan(
				"event_processed",
				opentracing.ChildOf(parentSpanCtx),
				opentracing.Tag{Key: "EventId", Value: eventStreamQuantum.GetEvent().GetEventId()},
			)
		}

		er.reducer.Reduce(eventStreamQuantum.GetEvent())

		if nil != eventSpan {
			eventSpan.Finish()
		}
	}

	if nil != streamClientSpan {
		streamClientSpan.Finish()
	}
}

func (er *EventReaderReducer) SubscribeReduceGlobalStream() {
	streamObj, streamError := er.client.SubscribeGlobalStream(context.Background(), &apiEvent.SubscribeGlobalStreamRequest{})
	if nil != streamError {
		panic("failed to call GlobalStream: " + streamError.Error())
	}

	streamClientSpan := projectorOpentracing.SpanFromContext(streamObj.Context())

	for {
		eventStreamQuantum, err := streamObj.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}

			panic("error while reading global event stream: " + err.Error())
		}

		var eventSpan opentracing.Span
		if parentSpanCtx := er.selectOpentracingContextForEventStreamQuantum(streamClientSpan, eventStreamQuantum); nil != parentSpanCtx {
			eventSpan = er.tracer.StartSpan(
				"subscription__event_processed",
				opentracing.ChildOf(parentSpanCtx),
				opentracing.Tag{Key: "EventId", Value: eventStreamQuantum.GetEvent().GetEventId()},
			)
		}

		er.reducer.Reduce(eventStreamQuantum.GetEvent())

		if nil != eventSpan {
			eventSpan.Finish()
		}
	}
}

func (er *EventReaderReducer) selectOpentracingContextForEventStreamQuantum(
	streamClientSpan opentracing.Span,
	eventStreamQuantum *apiEvent.EventStreamQuantum,
) opentracing.SpanContext {
	var selectedSpanCtx opentracing.SpanContext

	quantumOpentracingCtx, quantumOpentracingCtxError := er.tracer.Extract(
		opentracing.TextMap,
		opentracing.TextMapCarrier(eventStreamQuantum.GetMeta()),
	)
	if nil != quantumOpentracingCtxError {
		quantumOpentracingCtx = nil // To be sure that it is null
		fmt.Printf("Failed to extract Opentracing context from EventStreamQuantum \n")
	}

	if nil != quantumOpentracingCtx {
		selectedSpanCtx = quantumOpentracingCtx
	} else if nil != streamClientSpan {
		selectedSpanCtx = streamClientSpan.Context()
	}

	return selectedSpanCtx
}

func NewEventReader(client apiEvent.APIClient, otBridge *projectorOpentracing.Bridge, reducer *EventReducer) *EventReaderReducer {
	return &EventReaderReducer{
		client:  client,
		tracer:  otBridge.Tracer(),
		reducer: reducer,
	}
}
