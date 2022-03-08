package polytree

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/do87/poly/src/agent/internal/logger"
)

// Tree is the polytree handler
type Tree struct {
	Key     string
	RunID   string
	Nodes   []*Node
	Timeout time.Duration
	Meta    interface{}
	Payload []byte //request payload

	// internal
	seenNodes   map[string]bool
	errors      map[string]error
	pendingRun  []*Node
	pendingLock *sync.Mutex
}

// Node is a node in the polytree
type Node struct {
	Key      string
	Parents  []*Node
	Children []*Node
	Exec     Exec
	Error    error
}

type Exec func(ctx context.Context, log *logger.Logger, meta interface{}, payload []byte) (Exec, error)

func New() *Tree {
	return &Tree{
		Timeout: 1 * time.Hour,
	}
}

// AddNode adds a node to the polytree
func (t *Tree) AddNode(node *Node) *Tree {
	t.Nodes = append(t.Nodes, node)
	return t
}

// AddDependency creates dependencies between nodes
func (t *Tree) AddDependency(parent, child *Node) *Tree {
	parent.Children = append(parent.Children, child)
	child.Parents = append(child.Parents, parent)
	return t
}

func (t *Tree) ParentOf(parent, child *Node) *Tree {
	return t.AddDependency(parent, child)
}

func (t *Tree) execNode(ctx context.Context, node *Node, done chan *Node) {
	ctxWrap, cancel := context.WithTimeout(ctx, t.Timeout)
	if !t.shouldNodeRun(ctxWrap, cancel, node) {
		return
	}

	log, logd := logger.NewNodeLogger(node.Key, t.RunID)
	defer logd()

	var err error
	ch := make(chan error)
	go func() {
		// catch panic
		defer func() {
			if err := recover(); err != nil {
				ch <- fmt.Errorf("%v", err)
			}
		}()

		// run all steps
		for step, err := node.Exec(ctxWrap, log, t.Meta, t.Payload); step != nil; {
			if err != nil {
				ch <- err
				break
			}

			// cancel if context expired
			if err = ctxWrap.Err(); err != nil {
				ch <- err
				break
			}

			// run next step
			if step, err = step(ctxWrap, log, t.Meta, t.Payload); err != nil {
				ch <- err
			}
		}
		if err != nil {
			ch <- err
		}
		ch <- nil
	}()

	select {
	case err := <-ch:
		log.Error(err)
		node.Error = err
	case <-ctxWrap.Done():
		node.Error = errors.New("node execution timeout")
	}

	t.setSeen(node)
}

// shouldNodeRun returns true if the node should run
// if one of the parents hasn't run yet, the function will wait for it
func (t *Tree) shouldNodeRun(ctx context.Context, cancel context.CancelFunc, node *Node) bool {
	for {
		// don't run node if there are errors
		if len(t.errors) > 0 {
			cancel()
			node.Error = errors.New("execution skipped")
			t.setSeen(node)
			return false
		}

		// check context in case of timeout or sigterm
		if err := ctx.Err(); err != nil {
			cancel()
			node.Error = err
			t.setSeen(node)
			return false
		}

		// validate parent nodes
		ready := true
		for _, parent := range node.Parents {
			if !t.isSeen(parent) {
				ready = false
			}
		}
		if ready {
			break
		}
		time.Sleep(time.Second * 5)
	}
	return true
}

// Execute is used to execute each polytree node.Exec function
func (t *Tree) Execute(ctx context.Context, log *logger.Logger) {
	t.pendingRun = t.getTopNodes()
	t.execute(ctx)
	t.cleanup()
}

// ExecuteWithTimeout is used to execute each polytree node.Exec function, with a global run timeout
func (t *Tree) ExecuteWithTimeout(ctx context.Context, timeout time.Duration) {
	ctxWrap, cancel := context.WithTimeout(ctx, timeout)
	t.pendingRun = t.getTopNodes()
	t.execute(ctxWrap)
	t.cleanup(cancel)
}

// execute runs all pending nodes
func (t *Tree) execute(ctx context.Context) {
	if len(t.pendingRun) == 0 {
		return
	}

	// Get nodes that need to run, clear them from the pending list
	t.pendingLock.Lock()
	nodes := t.pendingRun
	t.removeFromPending(nodes...)
	t.pendingLock.Unlock()

	done := make(chan *Node, len(nodes))
	for _, node := range nodes {
		go t.execNode(ctx, node, done)
	}

	for i := 0; i < len(nodes); i++ {
		node := <-done
		t.setSeen(node)

		// add children to pending list
		t.pendingLock.Lock()
		t.pendingRun = append(t.pendingRun, node.Children...)
		t.pendingLock.Unlock()
		t.execute(ctx)
	}
}

// removeFromPending remove given nodes from pending list
func (t *Tree) removeFromPending(nodes ...*Node) {
	newPending := []*Node{}
PLOOP:
	for _, pending := range t.pendingRun {
		for _, n := range nodes {
			if pending == n {
				continue PLOOP
			}
		}
		newPending = append(newPending, pending)
	}
	t.pendingRun = newPending
}

// setSeen marks a node as completed / seen
func (t *Tree) setSeen(n *Node) {
	t.seenNodes[n.Key] = true
	if n.Error != nil {
		t.errors[n.Key] = n.Error
	}
}

// isSeen returns true if node completed
func (t *Tree) isSeen(n *Node) bool {
	v, ok := t.seenNodes[n.Key]
	return ok && v
}

// getTopNodes returns all nodes without a parent
func (t *Tree) getTopNodes() []*Node {
	n := []*Node{}
	for _, v := range t.Nodes {
		if len(v.Parents) == 0 {
			n = append(n, v)
		}
	}
	return n
}

func (t *Tree) cleanup(cancel ...context.CancelFunc) {
	t.seenNodes = map[string]bool{}
	t.errors = map[string]error{}
	t.pendingRun = []*Node{}
	for _, c := range cancel {
		c()
	}
}
