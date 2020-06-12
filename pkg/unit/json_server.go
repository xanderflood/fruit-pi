package unit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func init() {
	RegisterUnitType("json_server", Schema{
		GetEmptyConfigPointer: func() TypeConfig { return &JSONServerConfig{} },
		Inputs:                map[string]struct{}{},
		Outputs: map[string]OutputSchema{
			"triggered": {
				NoCaching: true,
			},
		},
	})
}

type JSONServerConfig struct {
	Port int `json:"port"`
}

func (cfg JSONServerConfig) New() UnitV2 {
	return &JSONServer{cfg: cfg}
}

type JSONServer struct {
	cfg JSONServerConfig
	srv *http.Server

	triggered Output
}

func (unit *JSONServer) Start(ctx context.Context, _ map[string]Input, outputs map[string]Output) <-chan struct{} {
	done := make(chan struct{}, 1)

	unit.srv = &http.Server{Addr: fmt.Sprintf(":%v", unit.cfg.Port)}
	unit.triggered = outputs["triggered"]

	go func() {
		defer close(done)

		if err := unit.srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %v", err)
		}
	}()

	return done
}

func (unit *JSONServer) handleRequest(w http.ResponseWriter, r *http.Request) {
	var val interface{}
	if err := json.NewDecoder(r.Body).Decode(&val); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	unit.triggered<-val

	return
}

func (unit *JSONServer) GetState() interface{} {
	return nil
}
func (unit *JSONServer) SetState(interface{}) {}
func (unit *JSONServer) Stop() {
	unit.srv.Shutdown(context.Background())
}
