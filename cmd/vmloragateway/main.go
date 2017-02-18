package main

import (
	//"encoding/base64"
	"os"
	"os/signal"
	"syscall"
	//"time"

	log "github.com/Sirupsen/logrus"
	"github.com/panyingyun/vmloragateway/backend"
	"github.com/panyingyun/vmloragateway/config"
	//gw "github.com/panyingyun/vmloragateway/gateway"
	"github.com/panyingyun/vmloragateway/server"
	_ "github.com/panyingyun/vmspace/gateway"
	_ "github.com/panyingyun/vmspace/node"
	_ "github.com/smallnest/rpcx"
	"github.com/urfave/cli"
)

func run(c *cli.Context) error {
	// Read Config
	conf, _ := config.ReadConfig(c.String("conf"))
	log.Info(conf)

	gatewayid := c.String("gateway")
	log.Info(gatewayid)

	// Connect to Lora-Gateway-Bridge
	backend, err := backend.NewBackend(conf.ServerAddr, gatewayid, conf.Longtitude, conf.Latitude, conf.Altitude)
	if err != nil {
		return err
	}
	// log.Infof("backend = %v, err = %v", backend, err)

	// Start Send Gateway Stat every 30s
	hbserver := server.NewHBServer(backend, conf, gatewayid)
	hbserver.Start()
	defer hbserver.Stop()

	//	// Connect to Rpc Server
	//	selector := &rpcx.DirectClientSelector{
	//		Network:     "tcp",
	//		Address:     "127.0.0.1:8972",
	//		DialTimeout: 10 * time.Second,
	//	}
	//	client := rpcx.NewClient(selector)
	//	defer client.Close()

	//	//try to get rx msg and tranmit to loarserver
	//	go func() {
	//		for {

	//			args4 := &gateway.GWReceiveArgs{
	//				Gwid: gatewayid,
	//			}
	//			var reply4 gateway.GWReceiveReply
	//			client.Call("GW.Receive", args4, &reply4)
	//			log.Println("GW.Receive = ", len(reply4.Payload))
	//			if len(reply4.Payload) > 0 {
	//				now := time.Now().UTC()
	//				rxpk := gw.RXPK{
	//					Time: CompactTime(now),
	//					Tmst: uint32(time.Now().UnixNano() / 1000000),
	//					Freq: 868.5,
	//					Chan: 2,
	//					RFCh: 1,
	//					Stat: 1,
	//					Modu: "LORA",
	//					DatR: DatR{LoRa: "SF7BW125"},
	//					CodR: "4/5",
	//					RSSI: -51,
	//					LSNR: 7,
	//					Size: 16,
	//					Data: base64.StdEncoding.EncodeToString(reply4.Payload),
	//				}
	//			}
	//			time.Sleep(time.Second)

	//		}
	//	}()

	//	//try to get tx msg from loarserver and tranmit to node
	//	//	go func() {
	//	//		for {
	//	//			//try to get received msg and tranmit to loarserver
	//	//			args4 := &gateway.GWReceiveArgs{
	//	//				Gwid: gatewayid,
	//	//			}
	//	//			var reply4 gateway.GWReceiveReply
	//	//			client.Call("GW.Receive", args4, &reply4)
	//	//			log.Println("GW.Receive = ", len(reply4.Payload))

	//	//			time.Sleep(time.Second)

	//	//		}
	//	//	}()

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
	app.Usage = "vmloragateway -gw F1E2D3C4B5A60000"
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
