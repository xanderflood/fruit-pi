package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

func Test1(t *testing.T) (string, error) {
	RegisterTestingT(t)

	server := ghttp.NewServer()
	defer server.Close()

	server.AppendHandlers(
		//TODO update this
		ghttp.CombineHandlers(
			ghttp.VerifyRequest("POST", "/v3/adgroup/"),
			ghttp.VerifyJSON(adgroupLifecycleSuccessJSONStep1(false, false)),
			ghttp.RespondWith(http.StatusInternalServerError, nil),
		),
	)

	docker, err := client.NewEnvClient()
	Expect(err).To(BeNil())

	//
	// TODO
	// 1) start the image with:
	// - host networking
	// - a config file passed in
	//
	//
	//
	//
	//
	//
	//
	//
	//

	cont, err := docker.ContainerCreate(context.Background(),
		&container.Config{
			Image: "xanderflood/fruit-pi",
			Cmd:   []string{"./controller", "--base64-device-config", "/static-config.json"},
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: fmt.Sprintf("%s/config.json", os.Getenv("PWD")),
					Target: "/static-config.json",
				},
			},
		},
		nil,
		"",
	)
	Expect(err).To(BeNil())

	docker.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	fmt.Printf("Container %s is started", cont.ID)
	return cont.ID, nil
}

func config(outputServerAddr string) string {
	device := Device{
		DeviceUUID: "device-uuid",
		Name:       "my special device name",
		Config:     &json.RawMessage([]byte(`

`)),
	}

	bs, err := json.Marshal(device)
	Expect(err).To(BeNil())

	return base64.RawStdEncoding.EncodeToString(bs)

}
