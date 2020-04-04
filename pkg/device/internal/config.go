package internal

import (
	"fmt"
	"strings"

	"github.com/twmb/algoimpl/go/graph"
	"github.com/xanderflood/fruit-pi/pkg/unit"
)

type Config struct {
	Version string                 `json:"version"`
	UUID    string                 `json:"uuid"`
	Units   map[string]unit.Config `json:"units"`
}

func (cfg Config) buildGraph() (Graph, error) {
	nodes := map[string]graph.Node{}
	outgoingEdges := map[string]edge{}

	g := graph.New(graph.Directed)
	for uName, u := range cfg.Units {
		node := g.MakeNode()

		var uI interface{} = &Node{
			name: uName,
			cfg:  u,
		}
		node.Value = &uI
		nodes[uName] = node

		for field, source := range u.Inputs {
			g.MakeEdge(nodes[source.Unit], node)

			if len(outgoingEdges[source.Unit]) == 0 {
				outgoingEdges[source.Unit] = edge{}
			}
			outgoingEdges[source.Unit][source.Name] = connection{
				source: source,
				dest: unit.ValueIdentifier{
					Unit: uName,
					Name: field,
				},
			}
		}
	}

	// check for cycles
	components := g.StronglyConnectedComponents()
	if len(nodes) != len(components) {
		for _, component := range components {
			if len(component) > 1 {
				names := make([]string, len(component))
				for i := range component {
					names[i] = (*component[i].Value).(*Node).name
				}
				return Graph{}, fmt.Errorf("circular reference:\n%s", strings.Join(names, "\n"))
			}
		}

		return Graph{}, fmt.Errorf("INACCESSIBLE POINT")
	}

	return Graph{Graph: g}, nil
}
