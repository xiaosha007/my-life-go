package monitoring

import (
	"fmt"
	"log"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/xiaosha007/my-life-go/pkg/transformer"
)

var StatsdClient *statsd.Client

type StatsdClientConfig struct {
	Namespace string `json:"namespace"`
	Host      string `json:"host"`
	Port      int    `json:"port"`
}

func NewStatsdClient(config *StatsdClientConfig) *statsd.Client {

	dogstatsd_client, err := statsd.New(fmt.Sprintf("%s:%s", config.Host, transformer.IntToString(config.Port)), statsd.WithNamespace(config.Namespace))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Created statsd client!")

	StatsdClient = dogstatsd_client

	return StatsdClient
}
