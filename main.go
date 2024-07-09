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

	fmt.Println("Subscribed to topic", "optimizations")
	_, err = c.Subscribe("optimizations", func(p *utils.OptimizationRequest) {
		splits := strings.Split(p.ResourceId, "/")
		name := splits[len(splits)-1] + "-opt"
		fmt.Println("Received optimization, publishing CR object", name)
		err := utils.CreateOptimizationCustomResource(p, name)
		utils.Fatal(err)
	})
	if err != nil {
		utils.Fatal(err)
	}

	for {
		time.Sleep(1 * time.Second)
	}
}
