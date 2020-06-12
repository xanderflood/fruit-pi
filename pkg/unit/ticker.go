package unit

import (
	"context"
	"time"

	"github.com/xanderflood/fruit-pi/lib/config"
)

func init() {
	RegisterUnitType("ticker", Schema{
		GetEmptyConfigPointer: func() TypeConfig { return &TickerConfig{} },
		Inputs:                map[string]struct{}{},
		Outputs: map[string]OutputSchema{
			"tick": {
				NoCaching: true,
			},
		},
	})
}

type TickerConfig struct {
	Interval config.Duration `json:"interval"`
}

func (tc TickerConfig) New() UnitV2 {
	return &Ticker{cfg: tc}
}

type Ticker struct {
	cfg    TickerConfig
	cancel context.CancelFunc
}

func (tu *Ticker) Start(ctx context.Context, _ map[string]Input, outputs map[string]Output) <-chan struct{} {
	done := make(chan struct{}, 1)

	cCtx, cancel := context.WithCancel(ctx)
	t := time.NewTicker(tu.cfg.Interval.Duration)
	tu.cancel = cancel
	go func() {
		defer close(done)

		for {
			select {
			case <-cCtx.Done():
				break
			case <-t.C:
				outputs["tick"] <- nil
			}
		}
	}()

	return done
}

func (t *Ticker) GetState() interface{} {
	return nil
}
func (t *Ticker) SetState(interface{}) {}
func (t *Ticker) Stop() {
	t.cancel()
}
