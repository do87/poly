package logger

import (
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/chizap"
)

// Logger service
type Logger struct {
	comp *zap.Logger
}

// New creates a new Logger
func New() (*Logger, func() error) {
	log, _ := zap.NewProduction()
	return &Logger{
		comp: log,
	}, log.Sync
}

// ChiMiddleware returns a chi specific logging middleware
func (l *Logger) ChiMiddleware() func(http.Handler) http.Handler {
	return chizap.New(l.comp, &chizap.Opts{
		WithReferer:   true,
		WithUserAgent: true,
	})
}

// NodeLogger is logger used during node run
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

// Named for namespacing logs
func (l *Logger) Named(name string) *Logger {
	return &Logger{
		comp: l.comp.Named(name),
	}
}

// Info log type
// every 2 args are added as key=value
func (l *Logger) Info(msg string, args ...interface{}) {
	l.comp.WithOptions(zap.AddCallerSkip(1)).Sugar().Infow(msg, args...)
}

// Warning log trype
// every 2 args are added as key=value
func (l *Logger) Warning(msg string, args ...interface{}) {
	l.comp.WithOptions(zap.AddCallerSkip(1)).Sugar().Warnw(msg, args...)
}

// Error log type
// every 2 args are added as key=value
func (l *Logger) Error(msg string, args ...interface{}) {
	l.comp.WithOptions(zap.AddCallerSkip(1)).Sugar().Errorw(msg, args...)
}

// Debug log type
// every 2 args are added as key=value
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.comp.WithOptions(zap.AddCallerSkip(1)).Sugar().Debugw(msg, args...)
}

// Fatal log type
// every 2 args are added as key=value
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.comp.WithOptions(zap.AddCallerSkip(1)).Sugar().Fatalw(msg, args...)
}
