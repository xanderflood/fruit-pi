package tools

import (
	"github.com/xanderflood/fruit-pi-server/lib/api"
)

//APIClient is used for generating fakes
//go:generate counterfeiter . APIClient
type APIClient interface {
	api.API
}
