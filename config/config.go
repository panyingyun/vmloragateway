package config

import (
	"fmt"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/ini.v1"
)

type Config struct {
	ServerAddr string  `ini:"server_addr"`
	GatewayMAC string  `ini:"gateway_id"`
	Latitude   float64 `ini:"latitude"`
	Longtitude float64 `ini:"longtitude"`
	Altitude   int32   `ini:"altitude"`

	LogDirWin   string `ini:"log_dir_win"`
	LogDirLinux string `ini:"log_dir_linux"`
	LogPrefix   string `ini:"log_prefix"`
}

func (c Config) String() string {
	server := fmt.Sprintf("Gateway:[%v]/[%v]/[%v]/[%v]", c.ServerAddr, c.GatewayMAC, c.Latitude, c.Longtitude, c.Altitude)

	log := fmt.Sprintf("LOG:[win:%v]/[linux:%v]:[prefix:%v]", c.LogDirWin, c.LogDirLinux, c.LogPrefix)

	return server + ", " + log
}

//Read Server's Config Value from "path"
func ReadConfig(path string) (Config, error) {
	var config Config
	conf, err := ini.Load(path)
	if err != nil {
		log.Println("load config file fail!")
		return config, err
	}
	conf.BlockMode = false
	err = conf.MapTo(&config)
	if err != nil {
		log.Println("mapto config file fail!")
		return config, err
	}
	return config, nil
}
