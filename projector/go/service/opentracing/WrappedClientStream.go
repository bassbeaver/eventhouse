package opentracing

import (
	"context"
	"google.golang.org/grpc"
)

type wrappedClientStream struct {
	grpc.ClientStream
	WrappedContext context.Context // WrappedContext is the wrapper's own Context. You can assign it.
}

// Context returns the wrapper's WrappedContext, overwriting the nested grpc.ServerStream.Context()
func (w *wrappedClientStream) Context() context.Context {
	return w.WrappedContext
}

// wrapClientStream returns a ClientStream that has the ability to overwrite context.
func wrapClientStream(stream grpc.ClientStream) *wrappedClientStream {
	if existing, ok := stream.(*wrappedClientStream); ok {
		return existing
	}
	return &wrappedClientStream{ClientStream: stream, WrappedContext: stream.Context()}
}
