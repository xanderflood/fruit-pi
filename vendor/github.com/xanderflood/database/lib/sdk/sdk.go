package dbsdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"

	perrors "github.com/pkg/errors"

	"github.com/xanderflood/database/lib/web"
)

//SDK minimal interface
type SDK interface {
	CreateTable(name string) error
	Insert(tableName string, document map[string]string) error
	Index(tableName string)
}

//Client standard SDK implementation
type Client struct {
	Origin *url.URL
	Client *http.Client
}

func New(originString string) (*Client, error) {
	origin, err := url.Parse(originString)
	if err != nil {
		return nil, perrors.Wrapf(err, "invalid origin string `%s`", originString)
	}

	return &Client{
		Origin: origin,
		Client: http.DefaultClient,
	}, nil
}

func (c *Client) CreateTable(name string) error {
	response := web.JSONStandardResponse{}

	body := struct {
		TableName string `json:"table_name"`
	}{TableName: name}

	//TODO use the routes package here
	code, err := c.do("POST", "/v1/createTable", body, &response)
	if err != nil {
		return err
	}

	if code != 200 {
		return fmt.Errorf("response code %v from CreateTable request: %s",
			code, response.Message)
	}

	return nil
}

func (c *Client) Insert(tableName string, document map[string]string) error {
	response := web.JSONStandardResponse{}

	//TODO use the routes package here
	code, err := c.do("POST",
		path.Join("/v1/insert/", tableName),
		document, &response)
	if err != nil {
		return err
	}

	if code != 200 {
		return fmt.Errorf("response code %v from CreateTable request: %s",
			code, response.Message)
	}

	return nil
}

func (c *Client) Index(tableName string) {
	//TODO
}

func (c *Client) do(method, path string, body interface{}, responseObj interface{}) (int, error) {
	endpoint, err := url.Parse(path)
	if err != nil {
		return 0, perrors.Wrapf(err, "invalid request path `%s`", path)
	}

	requestPath := c.Origin.ResolveReference(endpoint)

	bodydata, err := json.Marshal(body)
	if err != nil {
		return 0, perrors.Wrapf(err, "could not marshal request body `%v`", body)
	}

	bodybuf := bytes.NewReader(bodydata)
	req, err := http.NewRequest(method, requestPath.String(), bodybuf)
	if err != nil {
		return 0, perrors.Wrap(err, "failed to build request")
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return 0, perrors.Wrapf(err, "%s request failed", method)
	}

	code := resp.StatusCode

	respbuf := bytes.NewBuffer([]byte{})
	_, err = io.Copy(respbuf, resp.Body)
	if err != nil {
		return code, perrors.Wrap(err, "failed reading request body")
	}

	err = json.Unmarshal(respbuf.Bytes(), responseObj)
	if err != nil {
		return code, perrors.Wrapf(err, "failed unmarshalling response body `%s`", respbuf.String())
	}

	return code, nil
}
