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

func (l *Logger) Info(args ...interface{}) {
	l.comp.Sugar().Info(args)
}

func (l *Logger) Warning(args ...interface{}) {
	l.comp.Sugar().Warn(args)
}

func (l *Logger) Error(args ...interface{}) {
	l.comp.Sugar().Error(args)
}

func (l *Logger) Debug(args ...interface{}) {
	l.comp.Sugar().Debug(args)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.comp.Sugar().Fatal(args)
}
