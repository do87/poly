package logger

import (
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/chizap"
)

type Logger struct {
	comp *zap.Logger
}

func New() (*Logger, func() error) {
	log, _ := zap.NewProduction()
	return &Logger{
		comp: log,
	}, log.Sync
}

func (l *Logger) ChiMiddleware() func(http.Handler) http.Handler {
	return chizap.New(l.comp, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	})
}

func (l *Logger) NodeLogger(plan, runID, node string) *Logger {
	log := l.comp.WithOptions(zap.Fields(
		zap.Field{
			Key:    "node",
			Type:   zapcore.StringType,
			String: node,
		},
		zap.Field{
			Key:    "plan",
			Type:   zapcore.StringType,
			String: plan,
		},
	))

	return &Logger{
		comp: log.Named(runID),
	}
}

func (l *Logger) Named(name string) *Logger {
	return &Logger{
		comp: l.comp.Named(name),
	}
}

func (l *Logger) Info(msg string, args ...interface{}) {
	l.comp.WithOptions(zap.AddCallerSkip(1)).Sugar().Infow(msg, args...)
}

func (l *Logger) Warning(msg string, args ...interface{}) {
	l.comp.WithOptions(zap.AddCallerSkip(1)).Sugar().Warnw(msg, args...)
}

func (l *Logger) Error(msg string, args ...interface{}) {
	l.comp.WithOptions(zap.AddCallerSkip(1)).Sugar().Errorw(msg, args...)
}

func (l *Logger) Debug(msg string, args ...interface{}) {
	l.comp.WithOptions(zap.AddCallerSkip(1)).Sugar().Debugw(msg, args...)
}

func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.comp.WithOptions(zap.AddCallerSkip(1)).Sugar().Fatalw(msg, args...)
}
