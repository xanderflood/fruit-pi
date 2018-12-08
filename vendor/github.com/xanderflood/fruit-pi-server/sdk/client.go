package sdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

//go:generate counterfeiter . Client
type Client interface {
	Initialize(config []byte) (uuid string, err error)
	Update(uuid string, config []byte) error
	Fetch(uuid string) (config []byte, err error)
}

//Frontmatterer exposes a frontmatter handler
//go:generate counterfeiter . Frontmatterer
type Frontmatterer interface {
	Frontmatter(method, url string, body interface{}, response interface{}) (status int, err error)
}

type ClientAgent struct {
	Origin        string
	Token         string
	Client        *http.Client
	Frontmatterer Frontmatterer
}

func New(origin, token string) *ClientAgent {
	httpClient := *http.DefaultClient
	ca := &ClientAgent{
		Origin: origin,
		Token:  token,
		Client: &httpClient,
	}
	ca.Frontmatterer = ca
	return ca
}

func (ca *ClientAgent) Create(config []byte) (uuid string, err error) {
	request := CreateRequest{Config: string(config)}
	response := CreateResponse{}
	status, err := ca.Frontmatterer.Frontmatter("POST", "/v1/create", &request, &response)
	if err != nil {
		return "", err
	}
	if status != http.StatusOK {
		//TODO: use the ErrorResponse struct to get a message
		return "", fmt.Errorf("invalid response code: %v", status)
	}

	return response.UUID, nil
}

func (ca *ClientAgent) Configure(uuid string, config []byte) error {
	request := ConfigureRequest{UUID: uuid, Config: string(config)}
	status, err := ca.Frontmatterer.Frontmatter("PATCH", "/v1/configure", request, nil)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		//TODO: use the ErrorResponse struct to get a message
		return fmt.Errorf("invalid response code: %v", status)
	}

	return nil
}

func (ca *ClientAgent) Configuration(uuid string) (config []byte, err error) {
	url := fmt.Sprintf("/v1/configuration/%s", uuid)
	response := ConfigurationResponse{}
	status, err := ca.Frontmatterer.Frontmatter("GET", url, nil, &response)
	if status == http.StatusNotModified {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, fmt.Errorf("invalid response code: %v", status)
	}

	return []byte(response.Config), nil
}

func (ca *ClientAgent) Frontmatter(method, url string, body interface{}, response interface{}) (status int, err error) {
	var bodyReader io.Reader
	if body != nil {
		bodyData, err := json.Marshal(body)
		if err != nil {
			return 0, err
		}
		bodyReader = bytes.NewBuffer(bodyData)
	}

	url = fmt.Sprintf("http://%s%s", ca.Origin, url)
	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", ca.Token))

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := ca.Client.Do(req)
	if err != nil {
		return 0, err
	}

	if response != nil {
		respData, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return resp.StatusCode, err
		}

		err = json.Unmarshal(respData, response)
		if err != nil {
			return resp.StatusCode, err
		}
	}

	return resp.StatusCode, nil
}
