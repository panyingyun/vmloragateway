package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/panyingyun/vmloragateway/backend"
	"github.com/panyingyun/vmloragateway/config"
	"github.com/panyingyun/vmloragateway/gateway"
	"github.com/panyingyun/vmloragateway/server"
	"github.com/urfave/cli"
)

func run(c *cli.Context) error {
	// Read Config
	conf, _ := config.ReadConfig(c.String("conf"))
	log.Info(conf)

	gatewayid := c.String("gateway")
	log.Info(gatewayid)

	// Connect to Lora-Gateway-Bridge
	backend, err := backend.NewBackend(conf.ServerAddr)
	log.Infof("backend = %v, err = %v", backend, err)

	// Start Send Gateway Stat every 30s
	log.Infof("PullACK = %v", gateway.PullACK)
	hbserver := server.NewHBServer(backend, conf, gatewayid)
	hbserver.Start()
	defer hbserver.Stop()

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
	app.Usage = "vmloragateway -gw F1E2D3C4B5A61314"
	app.Copyright = "panyingyun@gmail.com "
	app.Version = "0.1"
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "conf,c",
			Usage:  "Set conf path here",
			Value:  "gateway.conf",
			EnvVar: "VMGATEWAY_CONF",
		},
		cli.StringFlag{
			Name:   "gateway,gw",
			Usage:  "Set gateway ID here",
			Value:  "F1E2D3C4B5A60000",
			EnvVar: "GATEWAY_ID",
		},
	}
	app.Run(os.Args)
}
