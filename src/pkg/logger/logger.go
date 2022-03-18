package logger

import (
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/chizap"
)

// Log is the interface used by the agent and mesh server
type Log interface {
	ChiMiddleware() func(http.Handler) http.Handler
	NodeLogger(plan, runID, node string) Log
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Warning(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
}

var _ = Log(&Logger{})

// Logger service
type Logger struct {
	comp *zap.Logger
}

// New creates a new Logger
func New(opts ...zap.Option) (*Logger, func() error) {
	log, _ := zap.NewProduction(opts...)
	return &Logger{
		comp: log,
	}, log.Sync
}

// NewDevelopment creates a new Logger for development
func NewDevelopment(opts ...zap.Option) (*Logger, func() error) {
	log, _ := zap.NewDevelopment(opts...)
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
func (l *Logger) NodeLogger(plan, runID, node string) Log {
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

	return Log(&Logger{
		comp: log.Named(runID),
	})
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
