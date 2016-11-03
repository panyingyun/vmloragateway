package server

import (
	"github.com/panyingyun/vmloragateway/backend"
	"github.com/panyingyun/vmloragateway/config"
	"github.com/robfig/cron"

	log "github.com/Sirupsen/logrus"
)

var spec_heartbeat = "@every 10s"
var spec_stat = "@every 30s"

type HBServer struct {
	backend *backend.Backend
	conf    config.Config
	gwid    string
	cron    *cron.Cron
}

func NewHBServer(backend *backend.Backend, conf config.Config, gwid string) *HBServer {
	c := cron.New()
	c.AddFunc(spec_heartbeat, func() {
		//log.Println("heartbeat")
		err := backend.SendHeartbeat()
		if err != nil {
			log.Warnf("SendHeartbeat err = %v", err)
		}
	})

	c.AddFunc(spec_stat, func() {
		//log.Println("stat")
		err := backend.SendStatData()
		if err != nil {
			log.Warnf("SendStatData err = %v", err)
		}
	})
	return &HBServer{
		backend: backend,
		conf:    conf,
		gwid:    gwid,
		cron:    c,
	}
}

func (s *HBServer) Start() {
	s.cron.Start()
	log.Info("heartbeat server start")
}

func (s *HBServer) Stop() {
	s.cron.Stop()
	log.Info("heartbeat server stop")
}
