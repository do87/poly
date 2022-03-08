package logger

import "go.uber.org/zap"

type Logger struct {
	comp *zap.Logger
}

func New() (*Logger, func() error) {
	log, _ := zap.NewProduction()
	return &Logger{
		comp: log,
	}, log.Sync
}

func NewNodeLogger(name, run string) (*Logger, func() error) {
	log, _ := zap.NewProduction(zap.Fields(
		zap.Field{
			Key:    "run",
			String: run,
		},
	))

	return &Logger{
		comp: log.Named(name),
	}, log.Sync
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
