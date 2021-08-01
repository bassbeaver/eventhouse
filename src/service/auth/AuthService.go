package auth

import (
	"context"
	"database/sql"
	"encoding/base64"
	loggerService "github.com/bassbeaver/eventhouse/service/logger"
	opentracingBridge "github.com/bassbeaver/eventhouse/service/opentracing"
	"github.com/bassbeaver/logopher"
	"github.com/opentracing/opentracing-go"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

const (
	AuthServiceAlias = "AuthService"
	tokenTypeBasic   = "Basic"
)

type apiClient struct {
	ClientId   string
	SecretHash string
}

type AuthService struct {
	dbConnect         *sql.DB
	opentracingBridge *opentracingBridge.Bridge
	knownClients      map[string]*apiClient
}

func (s *AuthService) Auth(ctx context.Context) (context.Context, error) {
	if opentracingRootSpan := opentracing.SpanFromContext(ctx); nil != opentracingRootSpan {
		childSpan := s.opentracingBridge.Tracer().StartSpan(
			"interceptor__auth",
			opentracing.ChildOf(opentracingRootSpan.Context()),
		)
		defer childSpan.Finish()
	}
	logger := loggerService.GetLoggerFromContext(ctx)

	metadataObj, metadataIsOk := metadata.FromIncomingContext(ctx)
	if !metadataIsOk {
		logger.Warning("Auth error, request metadata not found", nil)

		return ctx, status.Error(codes.Unauthenticated, "request metadata not found")
	}

	authHeaders, authHeaderIsOk := metadataObj["authorization"]
	if !authHeaderIsOk {
		logger.Warning("Auth error, auth header not found", nil)

		return ctx, status.Error(codes.Unauthenticated, "authorization header not found")
	}

	authHeaderParts := strings.SplitN(authHeaders[0], " ", 2)
	if tokenTypeBasic != authHeaderParts[0] {
		logger.Warning("Auth error, invalid token type", &logopher.MessageContext{"token_type": authHeaderParts[0]})

		return ctx, status.Error(codes.Unauthenticated, "invalid access token type")
	}

	decodedCreds, decodedCredsError := base64.StdEncoding.DecodeString(authHeaderParts[1])
	if nil != decodedCredsError {
		logger.Warning("Failed to base64 decode auth credentials", &logopher.MessageContext{"error": decodedCredsError.Error()})
	}

	creds := strings.SplitN(string(decodedCreds), ":", 2)

	apiClientObj, apiClientExists := s.knownClients[creds[0]]
	if !apiClientExists {
		return ctx, status.Error(codes.Unauthenticated, "")
	}

	secretCheckError := bcrypt.CompareHashAndPassword([]byte(apiClientObj.SecretHash), []byte(creds[1]))
	if nil != secretCheckError {
		return ctx, status.Error(codes.Unauthenticated, "")
	}

	return ctx, nil
}

func (s *AuthService) loadKnownClients() {
	rows, queryError := s.dbConnect.Query("SELECT ClientId, SecretHash FROM apiClients")
	if nil != queryError {
		panic("Failed to query API Client: " + queryError.Error())
	}

	s.knownClients = make(map[string]*apiClient)

	for rows.Next() {
		apiClientObj := &apiClient{}
		scanError := rows.Scan(&apiClientObj.ClientId, &apiClientObj.SecretHash)
		if nil != scanError {
			panic("Failed to scan API Client query results: " + scanError.Error())
		}

		s.knownClients[apiClientObj.ClientId] = apiClientObj
	}
}

// --------

func NewAuthService(dbConnect *sql.DB, opentracingBridge *opentracingBridge.Bridge) *AuthService {
	s := &AuthService{dbConnect: dbConnect, opentracingBridge: opentracingBridge}
	s.loadKnownClients()

	return s
}
