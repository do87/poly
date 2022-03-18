package definitions

import (
	"reflect"
	"testing"

	"github.com/do87/poly/src/mesh/models"
)

func Test_assignmentProcess(t *testing.T) {
	type args struct {
		runs   []models.Run
		agents []models.Agent
	}
	type test struct {
		name string
		args args
		want map[string]string
	}
	tests := []test{
		{
			"test-1",
			args{
				runs:   []models.Run{{UUID: "run-1", Plan: "plan:infra:v1", Labels: []string{"prod", "infra"}}},
				agents: []models.Agent{{UUID: "agent-1", Active: true, Plans: []string{"plan:infra:v1"}, Labels: []string{"infra", "prod"}}},
			},
			map[string]string{"run-1": "agent-1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := assignmentProcess(tt.args.runs, tt.args.agents)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("assignmentProcess() = %v, want %v", got, tt.want)
			}
		})
	}
}
