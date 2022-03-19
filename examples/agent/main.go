package main

import (
	"context"
	"os"

	"github.com/do87/poly/example/agent/plans/infra"
	"github.com/do87/poly/src/agent"
	"github.com/do87/poly/src/pkg/auth"
	"github.com/do87/poly/src/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	ctx := context.Background()
	log, logsync := logger.NewCustom(getCustomLogger())
	defer logsync()

	agent.New(agent.Config{
		Logger: log,
		Labels: agent.Labels{
			"infra", "prod",
		},
		Key: auth.AgentRegisterKey{
			Name:       "pubkey:v0.1",
			PrivateKey: getPrivateKey(),
		},
		AgentHost: "localhost",
		MeshURL:   "http://127.0.0.1:8080",
	}).Plans(
		infra.Plan(),
	).Run(ctx)
}

func getPrivateKey() []byte {
	k, err := os.ReadFile("keys/private_key.pem")
	if err != nil {
		panic(err)
	}
	return k
}

func getCustomLogger() *zap.Logger {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.StacktraceKey = zapcore.OmitKey
	cfg.EncoderConfig.CallerKey = zapcore.OmitKey
	cfg.EncoderConfig.TimeKey = zapcore.OmitKey
	cfg.EncoderConfig.NameKey = zapcore.OmitKey
	cfg.Level.SetLevel(zapcore.InfoLevel)
	custom, _ := cfg.Build()
	return custom
}
