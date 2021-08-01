package api_client

import (
	"context"
	"fmt"
	"github.com/bassbeaver/eventhouse/api/compiled/event"
	"github.com/bassbeaver/eventhouse/projector/service/opentracing"
	"google.golang.org/grpc"
	"time"
)

const (
	ApiClientServiceAlias = "ApiClientService"
)

type ApiClient struct {
	event.APIClient
	connection *grpc.ClientConn
}

func (c *ApiClient) CloseConnection() error {
	return c.connection.Close()
}

func NewApiClient(
	apiHost string,
	apiPort int,
	authToken string,
	opentracingBridge *opentracing.Bridge,
) *ApiClient {
	authInterceptor := newClientAuthInterceptor(authToken)

	connectionCtx, connectionCtxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer connectionCtxCancel()

	apiConnection, apiConnectionError := grpc.DialContext(
		connectionCtx,
		fmt.Sprintf("%s:%d", apiHost, apiPort),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithChainUnaryInterceptor(
			opentracingBridge.Intercept,
			authInterceptor.Intercept,
		),
		grpc.WithChainStreamInterceptor(
			opentracingBridge.InterceptStream,
			authInterceptor.InterceptStream,
		),
	)
	if apiConnectionError != nil {
		panic("failed to dial gRPC server: " + apiConnectionError.Error())
	}

	return &ApiClient{
		APIClient:  event.NewAPIClient(apiConnection),
		connection: apiConnection,
	}
}
