package api_client

import (
	"context"
	"google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"
)

type clientAuthInterceptor struct {
	authToken string
}

func (i *clientAuthInterceptor) Intercept(
	ctx context.Context,
	method string,
	req,
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return invoker(
		grpcMetadata.NewOutgoingContext(ctx, grpcMetadata.Pairs("Authorization", "Basic "+i.authToken)),
		method,
		req,
		reply,
		cc,
		opts...,
	)
}

func (i *clientAuthInterceptor) InterceptStream(
	ctx context.Context,
	desc *grpc.StreamDesc,
	cc *grpc.ClientConn,
	method string,
	streamer grpc.Streamer,
	opts ...grpc.CallOption,
) (grpc.ClientStream, error) {
	return streamer(
		grpcMetadata.AppendToOutgoingContext(ctx, "Authorization", "Basic "+i.authToken),
		desc,
		cc,
		method,
		opts...,
	)
}

func newClientAuthInterceptor(authToken string) *clientAuthInterceptor {
	return &clientAuthInterceptor{authToken: authToken}
}
