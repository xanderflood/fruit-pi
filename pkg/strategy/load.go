package strategy

import "github.com/xanderflood/fruit-pi/pkg/chamber"

func Load(name string) chamber.Strategy {
	switch name {
	case "standard":
		return &Standard{}
	default:
		return nil
	}
}
