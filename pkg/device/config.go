package device

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

func (cfg Config) BuildGraph() (Graph, error) {
	nodes := map[string]graph.Node{}
	outgoingEdges := map[string]edge{}

	g := graph.New(graph.Directed)
	for uName, u := range cfg.Units {
		node := g.MakeNode()

		*node.Value = newNode(uName, u)
		nodes[uName] = node
	}

	for uName, u := range cfg.Units {
		node := nodes[uName]
		for field, source := range u.Inputs {
			if err := g.MakeEdge(nodes[source.Unit], node); err != nil {
				// should be inaccessible
				return Graph{}, err
			}

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

	if err := cfg.assertAcyclic(g, len(nodes)); err != nil {
		return Graph{}, err
	}

	return Graph{Graph: g}, nil
}

func (cfg Config) assertAcyclic(g *graph.Graph, numNodes int) error {
	components := g.StronglyConnectedComponents()
	if numNodes != len(components) {
		for _, component := range components {
			if len(component) > 1 {
				names := make([]string, len(component))
				for i := range component {
					names[i] = (*component[i].Value).(*Node).name
				}
				return fmt.Errorf("circular reference:\n%s", strings.Join(names, "\n"))
			}
		}

		return fmt.Errorf("INACCESSIBLE POINT")
	}

	return nil
}
