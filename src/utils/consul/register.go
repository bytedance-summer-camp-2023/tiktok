package consul

import (
	"fmt"
	capi "github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"strconv"
	"tiktok/src/constant/config"
	"tiktok/src/utils/logging"
)

var consulClient *capi.Client

func init() {
	cfg := capi.DefaultConfig()
	cfg.Address = config.EnvCfg.ConsulAddr
	if c, err := capi.NewClient(cfg); err == nil {
		consulClient = c
		return
	} else {
		logging.Logger.Errorf("Connect Consul happens error: %v", err)
	}
}

func RegisterConsul(name string, port string) error {
	parsedPort, err := strconv.Atoi(port[1:]) // port start with ':' which like ':37001'
	logging.Logger.WithFields(log.Fields{
		"name": name,
		"port": parsedPort,
	}).Infof("Services Register Consul")

	if err != nil {
		return err
	}
	reg := &capi.AgentServiceRegistration{
		ID:   fmt.Sprintf("%s-1", name),
		Name: name,
		Port: parsedPort,
		Check: &capi.AgentServiceCheck{
			Interval:                       "5s",
			Timeout:                        "5s",
			GRPC:                           fmt.Sprintf("%s:%d/Heath", "192.168.31.110", parsedPort),
			DeregisterCriticalServiceAfter: "30s",
		},
	}
	if err := consulClient.Agent().ServiceRegister(reg); err != nil {
		return err
	}
	return nil
}
