package unit

import (
	"context"

	"github.com/xanderflood/fruit-pi/pkg/unit/statsd"
)

func init() {
	RegisterUnitType("statsd_count", Schema{
		GetEmptyConfigPointer: func() TypeConfig { return &CountStatConfig{} },
		Inputs: map[string]struct{}{
			"trigger": {},
		},
		Outputs: map[string]OutputSchema{},
	})
}

type CountStatConfig struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (cfg CountStatConfig) New() UnitV2 {
	return &CountStat{cfg: cfg}
}

type CountStat struct {
	cfg    CountStatConfig
	stats  statsd.StatsD
	cancel context.CancelFunc
}

func (unit *CountStat) Start(ctx context.Context, inputs map[string]Input, _ map[string]Output) <-chan struct{} {
	done := make(chan struct{}, 1)

	// TODO have a single shared client
	unit.stats = statsd.New(unit.cfg.Address)
	unit.stats.Start()

	cCtx, cancel := context.WithCancel(ctx)
	unit.cancel = cancel
	go func() {
		defer close(done)

		for {
			select {
			case <-cCtx.Done():
				break
			case <-inputs["trigger"]:
				unit.stats.Count(unit.cfg.Name, 1)
			}
		}
	}()

	return done
}

func (unit *CountStat) GetState() interface{} {
	return nil
}
func (unit *CountStat) SetState(interface{}) {}
func (unit *CountStat) Stop() {
	unit.cancel()
	unit.stats.Stop()
}
