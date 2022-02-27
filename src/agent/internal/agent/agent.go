package agent

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type agent struct {
	config Config
}

type Config struct {
}

func New(c Config) *agent {
	return &agent{
		config: c,
	}
}

func (a *agent) Run() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	ticker := time.NewTicker(1 * time.Second)

	select {
	case <-ctx.Done():
		a.eol(stop)
		return
	case <-ticker.C:
		return
	}
}
