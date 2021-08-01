package recovery

import (
	"context"
	"fmt"
	loggerService "github.com/bassbeaver/eventhouse/service/logger"
	"github.com/bassbeaver/logopher"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
)

const (
	RecoveryServiceAlias = "RecoveryService"
)

type RecoveryService struct {
	loggerObj *logopher.Logger
}

func (rs *RecoveryService) Recover(contextObj context.Context, recoveredError interface{}) (err error) {
	loggerObj := loggerService.GetLoggerOrNilFromContext(contextObj)
	if nil == loggerObj {
		loggerObj = rs.loggerObj
	}

	loggerObj.Critical(
		"Паника",
		&logopher.MessageContext{"error": fmt.Sprintf("%+v", recoveredError), "stack": string(debug.Stack())},
	)

	return status.Error(codes.Internal, "")
}

func NewRecoveryService(loggerFactory *loggerService.LoggerFactory) *RecoveryService {
	loggerObj, _ := loggerFactory.CreateLogger(1)

	return &RecoveryService{
		loggerObj: loggerObj,
	}
}
