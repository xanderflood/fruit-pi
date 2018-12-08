package sdk

import (
	"errors"

	"github.com/xanderflood/fruit-pi-server/lib/tools"
)

//Request validatable interface request
//go:generate counterfeiter . Request
type Request interface {
	Validate() error
}

//CreateRequest CreateRequest
type CreateRequest struct {
	Config string `json:"config"`
}

func (r *CreateRequest) Validate() error {
	//TODO

	return nil
}

//CreateResponse CreateRequest
type CreateResponse struct {
	UUID string `json:"uuid"`
}

//ConfigureRequest ConfigureRequest
type ConfigureRequest struct {
	UUID   string `json:"uuid"`
	Config string `json:"config"`
}

func (r *ConfigureRequest) Validate() error {
	if !tools.IsUUID(r.UUID) {
		return errors.New("malformed UUID")
	}

	return nil
}

//ConfigurationRequest ConfigurationRequest
type ConfigurationRequest struct {
	UUID string `json:"uuid"`
}

func (r *ConfigurationRequest) Validate() error {
	if !tools.IsUUID(r.UUID) {
		return errors.New("malformed UUID")
	}

	return nil
}

//ConfigurationResponse ConfigurationRequest
type ConfigurationResponse struct {
	Config string `json:"config"`
}

//ErrorResponse ErrorResponse
//TODO add:
//  Code ErrorCode `json:"code"`
//to communicate specific information to the frontend
type ErrorResponse struct {
	Message string `json:"message"`
}
