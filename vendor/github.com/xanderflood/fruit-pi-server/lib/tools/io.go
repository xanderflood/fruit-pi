package tools

import "io"

//go:generate counterfeiter . ReadCloser
type ReadCloser interface {
	io.ReadCloser
}
