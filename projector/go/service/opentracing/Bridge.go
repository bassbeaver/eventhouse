package opentracing

import (
	"context"
	"fmt"
	eventhouseOpentracing "github.com/bassbeaver/eventhouse/service/opentracing"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"
	"io"
)

const (
	OpentracingBridgeServiceAlias = "OpentracingBridgeService"
	spanContextKey                = "span"
	methodSubscribeGlobalStream   = "/eventhouse.grpc.event.API/SubscribeGlobalStream"
)

type Bridge struct {
	tracer       opentracing.Tracer
	tracerCloser io.Closer
}

// CloseTracer closes tracer object owned by Bridge
// TODO implement graceful shutdown and use Bridge.CloseTracer() inside it
func (br *Bridge) CloseTracer() {
	if nil != br.tracerCloser {
		closeError := br.tracerCloser.Close()
		if nil != closeError {
			fmt.Printf("Failed to close JaeggerTracer: %s \n", closeError.Error())
		}
	}
}

func (br *Bridge) Tracer() opentracing.Tracer {
	return br.tracer
}

func (br *Bridge) Intercept(
	ctx context.Context,
	method string,
	req,
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	clientSpan := br.tracer.StartSpan("call:" + method)
	defer clientSpan.Finish()
	metadata := extractOrCreateMdFromCtxForModification(ctx)
	injectError := br.tracer.Inject(
		clientSpan.Context(),
		eventhouseOpentracing.GrpcMetadata,
		eventhouseOpentracing.GrpcMetadataCarrier(metadata),
	)
	if injectError != nil {
		fmt.Printf("Injection of Opentracing headers failed: %s \n", injectError.Error())
	}

	return invoker(
		grpcMetadata.NewOutgoingContext(ctx, metadata),
		method,
		req,
		reply,
		cc,
		opts...,
	)
}

// InterceptStream Do not forget to manually close Opentracing span provided in context of ClientStream
func (br *Bridge) InterceptStream(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	clientSpan := br.tracer.StartSpan("call:" + method)
	metadata := extractOrCreateMdFromCtxForModification(ctx)
	injectError := br.tracer.Inject(
		clientSpan.Context(),
		eventhouseOpentracing.GrpcMetadata,
		eventhouseOpentracing.GrpcMetadataCarrier(metadata),
	)
	if injectError != nil {
		fmt.Printf("Injection of Opentracing headers failed: %s \n", injectError.Error())
	}

	// For long living "subscriptions" requests we recognize root span as "request beginning" mark and finish it immediately.
	if methodSubscribeGlobalStream == method {
		clientSpan.Finish()
	}

	clientStream, e := streamer(
		grpcMetadata.NewOutgoingContext(ctx, metadata),
		desc,
		cc,
		method,
		opts...,
	)

	newsClCtx := context.WithValue(clientStream.Context(), spanContextKey, clientSpan)

	ws := wrapClientStream(clientStream)
	ws.WrappedContext = newsClCtx

	return ws, e
}

func extractOrCreateMdFromCtxForModification(ctx context.Context) grpcMetadata.MD {
	if md, ok := grpcMetadata.FromOutgoingContext(ctx); !ok {
		return md.Copy()
	}

	return grpcMetadata.New(nil)
}

func SpanFromContext(ctx context.Context) opentracing.Span {
	if span, spanIsOk := ctx.Value(spanContextKey).(opentracing.Span); spanIsOk {
		return span
	}

	return nil
}

func NewBridge(agentHost string, agentPort int, serviceName string) *Bridge {
	jaeggerConfigObj := jaegerConfig.Configuration{
		ServiceName: serviceName,
		Reporter: &jaegerConfig.ReporterConfig{
			LocalAgentHostPort: fmt.Sprintf("%s:%d", agentHost, agentPort),
		},
		Sampler: &jaegerConfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
	}

	tracer, tracerCloser, tracerError := jaeggerConfigObj.NewTracer()
	if nil != tracerError {
		fmt.Printf("Failed to create JaeggerTracer: %s \n", tracerError.Error())

		tracer = &opentracing.NoopTracer{}
	}

	return &Bridge{
		tracer:       tracer,
		tracerCloser: tracerCloser,
	}
}
