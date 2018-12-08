package web

import "net/http"

//go:generate counterfeiter . RoundTripper
type RoundTripper interface {
	http.RoundTripper
}
