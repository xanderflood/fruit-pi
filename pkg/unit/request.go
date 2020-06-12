package unit

import (
	"context"
	"net/http"
)

func init() {
	RegisterUnitType("request", Schema{
		GetEmptyConfigPointer: func() TypeConfig { return &RequestConfig{} },
		Inputs: map[string]struct{}{
			"trigger": {},
		},
		Outputs: map[string]OutputSchema{},
	})
}

type RequestConfig struct {
	url string
}

func (cfg RequestConfig) New() UnitV2 {
	return &Request{
		cfg: cfg,
	}
}

type Request struct {
	cfg    RequestConfig
	cancel context.CancelFunc
}

func (unit *Request) Start(ctx context.Context, inputs map[string]Input, _ map[string]Output) <-chan struct{} {
	done := make(chan struct{}, 1)

	cCtx, cancel := context.WithCancel(ctx)
	unit.cancel = cancel
	go func() {
		defer close(done)

		for {
			select {
			case <-cCtx.Done():
				break
			case <-inputs["trigger"]:
				http.Get(unit.cfg.url)
			}
		}
	}()

	return done
}
func (unit *Request) GetState() interface{} {
	return nil
}
func (unit *Request) SetState(interface{}) {}
func (unit *Request) Stop() {
	unit.cancel()
}
