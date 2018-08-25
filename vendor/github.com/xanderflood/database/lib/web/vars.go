package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

//VarsGetter gets mux vars
//go:generate counterfeiter . VarsGetter
type VarsGetter interface {
	Get(r *http.Request) map[string]string
}

//MuxVars standard VarsGetter
type MuxVars struct{}

func (MuxVars) Get(r *http.Request) map[string]string {
	return mux.Vars(r)
}
