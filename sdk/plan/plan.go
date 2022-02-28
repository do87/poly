package plan

import "github.com/do87/poly/sdk/job"

type Plan struct {
	Type string
	Key  string
	Name string
	Jobs []*job.Job
}

func New() *Plan {
	return &Plan{}
}
