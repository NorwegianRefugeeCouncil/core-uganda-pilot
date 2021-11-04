package logging

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
)

type loggingField int

const (
	requestIdKey loggingField = iota
	sessionIdKey
	serverNameKey
	pidKey
	exeKey
	operationKey
	middlewareKey
	handlerKey
	subjectKey
	storeKey
	storeActionKey
	fieldsKey
)

var encoderCfg = zapcore.EncoderConfig{
	MessageKey:     "msg",
	LevelKey:       "level",
	NameKey:        "logger",
	EncodeLevel:    zapcore.LowercaseLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
}
var level = zap.NewAtomicLevel()
var core = zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), os.Stdout, level)
var logger = zap.New(core).With(zap.Int("pid", os.Getpid()))

func init() {
	level.SetLevel(zapcore.DebugLevel)
}

func WithRequestID(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, requestIdKey, requestId)
}
func WithSessionID(ctx context.Context, sessionId string) context.Context {
	return context.WithValue(ctx, sessionIdKey, sessionId)
}
func WithServerName(ctx context.Context, serverName string) context.Context {
	return context.WithValue(ctx, serverNameKey, serverName)
}
func WithHandlerName(ctx context.Context, handlerName string) context.Context {
	return context.WithValue(ctx, handlerKey, handlerName)
}
func WithOperation(ctx context.Context, routeName string) context.Context {
	return context.WithValue(ctx, operationKey, routeName)
}
func WithPID(ctx context.Context) context.Context {
	return context.WithValue(ctx, pidKey, os.Getpid())
}
func WithExe(ctx context.Context) context.Context {
	return context.WithValue(ctx, exeKey, path.Base(os.Args[0]))
}
func WithAuthnSubject(ctx context.Context, subject string) context.Context {
	return context.WithValue(ctx, subjectKey, subject)
}
func WithMiddleware(ctx context.Context, middleware string) context.Context {
	return context.WithValue(ctx, middlewareKey, middleware)
}
func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	if currentFields, ok := ctx.Value(fieldsKey).([]zap.Field); ok {
		fields = append(fields, currentFields...)
	}
	return context.WithValue(ctx, fieldsKey, fields)
}

func WithStore(ctx context.Context, storeName, actionName string, fields ...zap.Field) context.Context {
	ctx = context.WithValue(ctx, storeKey, storeName)
	ctx = context.WithValue(ctx, storeActionKey, actionName)
	if len(fields) > 0 {
		return WithFields(ctx, fields...)
	}
	return ctx
}

func SetLogLevel(l zapcore.Level) {
	level.SetLevel(l)
}

func NewStoreLogger(ctx context.Context, storeName, actionName string, fields ...zap.Field) (context.Context, *zap.Logger) {
	ctx = WithStore(ctx, storeName, actionName, fields...)
	l := NewLogger(ctx)
	return ctx, l
}

func NewLogger(ctx context.Context) *zap.Logger {
	newLogger := logger
	if requestId, ok := ctx.Value(requestIdKey).(string); ok {
		newLogger = newLogger.With(zap.String("request_id", requestId))
	}
	if sessionId, ok := ctx.Value(sessionIdKey).(string); ok {
		newLogger = newLogger.With(zap.String("session_id", sessionId))
	}
	if serverName, ok := ctx.Value(serverNameKey).(string); ok {
		newLogger = newLogger.With(zap.String("server_name", serverName))
	}
	if handlerName, ok := ctx.Value(operationKey).(string); ok {
		newLogger = newLogger.With(zap.String("operation", handlerName))
	}
	if middleware, ok := ctx.Value(middlewareKey).(string); ok {
		newLogger = newLogger.With(zap.String("middleware", middleware))
	}
	if pid, ok := ctx.Value(pidKey).(int); ok {
		newLogger = newLogger.With(zap.Int("pid", pid))
	}
	if exe, ok := ctx.Value(exeKey).(string); ok {
		newLogger = newLogger.With(zap.String("exe", exe))
	}
	if sub, ok := ctx.Value(subjectKey).(string); ok {
		newLogger = newLogger.With(zap.String("subject", sub))
	}
	if store, ok := ctx.Value(storeKey).(string); ok {
		newLogger = newLogger.With(zap.String("store", store))
	}
	if storeAction, ok := ctx.Value(storeActionKey).(string); ok {
		newLogger = newLogger.With(zap.String("store_action", storeAction))
	}
	if fields, ok := ctx.Value(fieldsKey).([]zap.Field); ok {
		newLogger = newLogger.With(fields...)
	}

	return newLogger
}
