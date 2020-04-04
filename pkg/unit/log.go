package unit

import (
	"context"
	"fmt"
	"time"
)

func init() {
	RegisterUnitType("log", Schema{
		GetEmptyConfigPointer: func() TypeConfig { return &LogConfig{} },
		Inputs: map[string]struct{}{
			"trigger": struct{}{},
		},
		Outputs: map[string]OutputSchema{},
	})
}

type LogConfig struct{}

func (LogConfig) New() UnitV2 {
	return &Log{}
}

type Log struct {
	cancel context.CancelFunc
}

func (lu *Log) Start(ctx context.Context, inputs map[string]Input, _ map[string]Output) <-chan struct{} {
	done := make(chan struct{}, 1)

	cCtx, cancel := context.WithCancel(ctx)
	lu.cancel = cancel
	go func() {
		defer close(done)

		for {
			select {
			case <-cCtx.Done():
				break
			case <-inputs["trigger"]:
				fmt.Println(time.Now().Format(time.RFC3339))
			}
		}
	}()

	return done
}
func (lu *Log) GetState() interface{} {
	return nil
}
func (lu *Log) SetState(interface{}) {}
func (lu *Log) Stop() {
	lu.cancel()
}
