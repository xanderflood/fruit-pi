package device

import (
	"context"

	"github.com/twmb/algoimpl/go/graph"
	"github.com/xanderflood/fruit-pi/pkg/unit"
)

type edge map[string]connection

type connection struct {
	source unit.ValueIdentifier
	dest   unit.ValueIdentifier
}

type Device interface {
	Start() error
	Refresh() error
}

type Graph struct {
	*graph.Graph
}

func (g Graph) Start(ctx context.Context) error {
	var broadcasts = map[string]unit.Broadcasts{}
	err := g.topologically(func(node *Node) error {
		return node.Build(ctx, broadcasts)
	})
	if err != nil {
		return err
	}

	return g.topologically(func(node *Node) error {
		node.Start(ctx)
		return nil
	})
}

func (g Graph) Stop(ctx context.Context) {
	_ = g.topologically(func(node *Node) error {
		node.Stop(ctx)
		return nil
	})
}

func (g Graph) topologically(f func(*Node) error) error {
	for _, gNode := range g.Graph.TopologicalSort() {
		err := f((*gNode.Value).(*Node))
		if err != nil {
			return err
		}
	}
	return nil
}

type Node struct {
	name string
	cfg  unit.Config

	broadcasts    unit.Broadcasts
	subscriptions unit.Subscriptions
	unit          unit.UnitV2

	done <-chan struct{}
}

func newNode(name string, config unit.Config) *Node {
	return &Node{
		name: name,
		cfg:  config,
	}
}

func (node *Node) Build(ctx context.Context, broadcastsSoFar map[string]unit.Broadcasts) (err error) {
	node.unit, node.broadcasts, node.subscriptions, err = node.cfg.Build(ctx, node.name, broadcastsSoFar)
	broadcastsSoFar[node.name] = node.broadcasts
	return
}

func (node *Node) Start(ctx context.Context) {
	node.done = node.unit.Start(ctx, node.subscriptions.Inputs(), node.broadcasts.Outputs())
}

func (node *Node) Stop(ctx context.Context) {
	node.unit.Stop()

	// wait for the unit to finish shutting down
	for {
		select {
		case <-ctx.Done():
			return
		case <-node.done:
			return
		}
	}
}
