package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/panyingyun/vmloragateway/backend"
	"github.com/panyingyun/vmloragateway/gateway"
	"github.com/urfave/cli"

	log "github.com/Sirupsen/logrus"
)

func run(c *cli.Context) error {
	addr := c.String("addr")

	backend, err := backend.NewBackend(addr)
	log.Infof("backend = %v, err = %v", backend, err)

	log.Infof("PullACK = %v", gateway.PullACK)
	//quit when receive end signal
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	log.Infof("signal received signal %v", <-sigChan)
	log.Warn("shutting down server")
	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "Visual Machine(Lora Gateway) connect to lora-gateway-bridge for test loar server benchmark or others"
	app.Usage = "vmloragateway --addr 0.0.0.0:1680 or vmloragateway -a 0.0.0.0:1680"
	app.Copyright = "panyingyun@gmail.com"
	app.Version = "0.1"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "addr,a",
			Usage:  "Set bridge listen address(such as 0.0.0.0:1680) here",
			Value:  "0.0.0.0:1680",
			EnvVar: "BRIDGE_ADDR",
		},
	}
	app.Run(os.Args)
}
