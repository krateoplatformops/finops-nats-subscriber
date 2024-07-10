package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/krateoplatformops/finops-nats-subscriber/pkg/utils"
)

var (
	c *nats.EncodedConn
)

func main() {
	subTopic := os.Getenv("SUB_TOPIC")
	optSecretName := os.Getenv("OPT_SECRET_NAME")
	optSecretNamespace := os.Getenv("OPT_SECRET_NAMESPACE")
	optNamespace := os.Getenv("OPT_NAMESPACE")

	fmt.Println("Connecting to NATS server in:", os.Getenv("NATS_SERVICE_HOST")+":"+os.Getenv("NATS_SERVICE_PORT"))
	nc, err := nats.Connect(os.Getenv("NATS_SERVICE_HOST") + ":" + os.Getenv("NATS_SERVICE_PORT"))
	if err != nil {
		utils.Fatal(err)
	}
	defer nc.Close()

	c, err = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		utils.Fatal(err)
	}
	defer c.Close()

	fmt.Println("Subscribed to topic", subTopic)
	_, err = c.Subscribe(subTopic, func(p *utils.OptimizationRequest) {
		splits := strings.Split(p.ResourceId, "/")
		optName := splits[len(splits)-1] + "-opt"
		fmt.Println("Received optimization, publishing CR object", optName)
		err := utils.CreateOptimizationCustomResource(p, optName, optNamespace, optSecretName, optSecretNamespace)
		utils.Fatal(err)
	})
	if err != nil {
		utils.Fatal(err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
