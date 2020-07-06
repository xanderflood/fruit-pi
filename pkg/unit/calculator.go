package unit

import (
	"context"

	"github.com/robertkrimen/otto"
)

func init() {
	RegisterUnitType("formula", Schema{
		GetEmptyConfigPointer: func() TypeConfig { return &FormulaConfig{} },
		Inputs:                map[string]struct{}{},
		Outputs: map[string]OutputSchema{
			"result": {},
		},
	})
}

type FormulaConfig struct {
	Formula string `json:"formula"`
}

func (cfg FormulaConfig) New() UnitV2 {
	return &Formula{cfg: cfg}
}

type Formula struct {
	cfg    FormulaConfig
	state  map[string]float64
	cancel context.CancelFunc
}

func (unit *Formula) Start(ctx context.Context, inputs map[string]Input, _ map[string]Output) <-chan struct{} {
	done := make(chan struct{}, 1)

	cCtx, cancel := context.WithCancel(ctx)
	unit.cancel = cancel

	unit.state = map[string]float64{}
	for vName, vInput := range inputs {
		go func(name string, input Input) {
			defer close(done) // TODO make this a waitgroup?

			for {
				select {
				case <-cCtx.Done():
					break
				case <-inputs["trigger"]:
					// TODO handle panics
					unit.state[name] = (<-input).(float64)

					unit.evaluate(unit.state)
				}
			}
		}(vName, vInput)
	}

	return done
}

func (unit *Formula) GetState() interface{} {
	return nil
}
func (unit *Formula) SetState(interface{}) {}
func (unit *Formula) Stop() {
	unit.cancel()
}

func (unit *Formula) evaluate(vars map[string]float64) (float64, error) {
	vm := otto.New()

	// TODO cancellation
	// vm.Interrupt = make(chan func(), 1) // The buffer prevents blocking
	// go func() {
	// 	time.Sleep(2 * time.Second) // Stop after two seconds
	// 	vm.Interrupt <- func() {
	// 		panic(halt)
	// 	}
	// }()

	for k, v := range vars {
		vm.Set(k, v)
	}

	val, err := vm.Run(unit.cfg.Formula)
	if err != nil {
		return 0, err
	}

	return val.ToFloat()
}
