package logger

import (
	"context"
	"github.com/bassbeaver/eventhouse/service/request_id_setter"
	"github.com/bassbeaver/logopher"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"
)

const (
	RequestContextLoggerSetterServiceAlias = "RequestContextLoggerSetter"
	requestContextKey                      = "logger"
	messageContextRequestIdKey             = "request_id"
	messageContextMethodKey                = "method"
	messageContextSourceKey                = "source"
	messageContextSourceValue              = "go-grpc"
	messageContextHostKey                  = "host"
	messageContextDownstreamRemoteIpKey    = "downstream_remote_ip"
	messageContextXForwardedForKey         = "x_forwarded_for"
	messageContextUserAgentKey             = "user_agent"
	messageContextRefererKey               = "referer"
	metadataUserAgentKey                   = "user-agent"
	metadataRefererKey                     = "referer"
	metadataAuthorityKey                   = ":authority"
	metadataDownstreamRemoteIpKey          = "downstream-remote-ip"
	metadataXForwardedForKey               = "x-forwarded-for"
)

type RequestContextLoggerSetter struct {
	loggerFactory *LoggerFactory
}

func (i *RequestContextLoggerSetter) Intercept(
	contextObj context.Context,
	requestObj interface{},
	serverInfo *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	loggerObj := i.createLoggerObj(contextObj, serverInfo.FullMethod)

	responseObj, responseError := handler(context.WithValue(contextObj, requestContextKey, loggerObj), requestObj)

	// Async export of buffered logs and close of stream used not to
	go func() {
		loggerObj.ExportBufferedMessages()
		loggerObj.CloseStreams()
	}()

	return responseObj, responseError
}

func (i *RequestContextLoggerSetter) InterceptStream(
	serviceObj interface{},
	serverStreamObj grpc.ServerStream,
	serverInfo *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) error {
	loggerObj := i.createLoggerObj(serverStreamObj.Context(), serverInfo.FullMethod)

	newCtx := context.WithValue(serverStreamObj.Context(), requestContextKey, loggerObj)
	wrappedServerStreamObj := grpc_middleware.WrapServerStream(serverStreamObj)
	wrappedServerStreamObj.WrappedContext = newCtx

	responseError := handler(serviceObj, wrappedServerStreamObj)

	// Async export of buffered logs and close of stream used not to
	go func() {
		loggerObj.ExportBufferedMessages()
		loggerObj.CloseStreams()
	}()

	return responseError
}

func (i *RequestContextLoggerSetter) createLoggerObj(
	contextObj context.Context,
	fullMethod string,
) *logopher.Logger {
	requestIdProcessor := func(message *logopher.Message) *logopher.Message {
		newMessage := message.Clone()
		newMessage.Context[messageContextRequestIdKey] = request_id_setter.GetIdFromContext(contextObj)
		newMessage.Context[messageContextSourceKey] = messageContextSourceValue
		newMessage.Context[messageContextMethodKey] = fullMethod
		if metadata, metadataIsOk := grpcMetadata.FromIncomingContext(contextObj); metadataIsOk {
			if metadataEntry := metadata.Get(metadataAuthorityKey); 0 < len(metadataEntry) {
				newMessage.Context[messageContextHostKey] = metadataEntry[0]
			}
			if metadataEntry := metadata.Get(metadataUserAgentKey); 0 < len(metadataEntry) {
				newMessage.Context[messageContextUserAgentKey] = metadataEntry[0]
			}
			if metadataEntry := metadata.Get(metadataRefererKey); 0 < len(metadataEntry) {
				newMessage.Context[messageContextRefererKey] = metadataEntry[0]
			}
			if metadataEntry := metadata.Get(metadataDownstreamRemoteIpKey); 0 < len(metadataEntry) {
				newMessage.Context[messageContextDownstreamRemoteIpKey] = metadataEntry[0]
			}
			if metadataEntry := metadata.Get(metadataXForwardedForKey); 0 < len(metadataEntry) {
				newMessage.Context[messageContextXForwardedForKey] = metadataEntry[0]
			}
		}

		return newMessage
	}

	loggerObj, createLoggerError := i.loggerFactory.CreateLogger(50)
	if nil != createLoggerError {
		panic("failed to create logger: " + createLoggerError.Error())
	}

	if nil != loggerObj.GetHandlers() {
		for _, handler := range loggerObj.GetHandlers() {
			if processorHolder, isProcessorHolder := handler.(logopher.ProcessorHolder); isProcessorHolder {
				processorHolder.AddProcessor(requestIdProcessor)
			}
		}
	}

	return loggerObj
}

func NewRequestContextLoggerInterceptor(loggerFactory *LoggerFactory) *RequestContextLoggerSetter {
	return &RequestContextLoggerSetter{
		loggerFactory: loggerFactory,
	}
}

//--------------------

func GetLoggerFromContext(contextObj context.Context) *logopher.Logger {
	return contextObj.Value(requestContextKey).(*logopher.Logger)
}

func GetLoggerOrNilFromContext(contextObj context.Context) *logopher.Logger {
	l := contextObj.Value(requestContextKey)
	if nil == l {
		return nil
	}

	return l.(*logopher.Logger)
}
