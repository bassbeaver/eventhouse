package request_id_setter

import (
	"context"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"
)

const (
	RequestIdSetterServiceAlias = "RequestIdSetter"
	requestIdContextKey         = "request_id"
	metadataRequestIdKey        = "x-request-id"
)

type RequestIdSetter struct {
}

func (i *RequestIdSetter) Intercept(
	contextObj context.Context,
	requestObj interface{},
	_ *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	return handler(context.WithValue(contextObj, requestIdContextKey, getRequestId(contextObj)), requestObj)
}

func (i *RequestIdSetter) InterceptStream(
	serviceObj interface{},
	serverStreamObj grpc.ServerStream,
	_ *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	newCtx := context.WithValue(serverStreamObj.Context(), requestIdContextKey, getRequestId(serverStreamObj.Context()))

	wrappedServerStreamObj := grpc_middleware.WrapServerStream(serverStreamObj)
	wrappedServerStreamObj.WrappedContext = newCtx

	return handler(serviceObj, wrappedServerStreamObj)
}

func getRequestId(contextObj context.Context) string {
	var requestId string

	if metadata, metadataIsOk := grpcMetadata.FromIncomingContext(contextObj); metadataIsOk {
		if requestIdMetadata := metadata.Get(metadataRequestIdKey); 0 < len(requestIdMetadata) {
			requestId = requestIdMetadata[0]
		}
	}

	if "" == requestId {
		if requestIdUuid, requestIdGenerationError := uuid.NewUUID(); nil == requestIdGenerationError {
			requestId = requestIdUuid.String()
		} else {
			requestId = "request_id_generation_error:" + requestIdGenerationError.Error()
		}
	}

	return requestId
}

func NewRequestIdSetter() *RequestIdSetter {
	return &RequestIdSetter{}
}

//--------------------

func GetIdFromContext(contextObj context.Context) string {
	return contextObj.Value(requestIdContextKey).(string)
}
