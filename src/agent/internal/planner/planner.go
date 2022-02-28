package planner

import "github.com/do87/poly/src/agent/internal/polytree"

type planner polytree.Tree
type plan polytree.Node
type Exec polytree.Exec

// New returns a new planner
func New() *planner {
	super := polytree.New()
	p := planner(*super)
	return &p
}

// AddPlan adds a plan to the polytree
func (p *planner) AddPlan(planType, planKey string, exec Exec) *plan {
	pl := &polytree.Node{
		Type: planType,
		Key:  planKey,
		Exec: polytree.Exec(exec),
	}
	p.Nodes = append(p.Nodes, pl)
	return (*plan)(pl)
}

// Dependency creates dependencies between plans
func (p *planner) Dependency(parent, child *plan) *planner {
	parent.Children = append(parent.Children, (*polytree.Node)(child))
	child.Parents = append(child.Parents, (*polytree.Node)(parent))
	return p
}
