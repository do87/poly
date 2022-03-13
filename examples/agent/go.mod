module github.com/do87/poly/example/agent

go 1.17

replace github.com/do87/poly/src => ../../src

require github.com/do87/poly/src v0.0.0-00010101000000-000000000000

require (
	github.com/go-chi/chi v1.5.4 // indirect
	go.uber.org/atomic v1.8.0 // indirect
	go.uber.org/multierr v1.7.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	moul.io/chizap v1.0.3 // indirect
)
