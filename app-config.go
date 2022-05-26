package kconfigurator

import (
	"fmt"
	"os"

	"git.kanosolution.net/kano/kaos"
	"github.com/eaciit/toolkit"
)

type AppConfig struct {
	Hosts       map[string]string
	Connections map[string]struct {
		Txt      string
		PoolSize int
	}
	EventServer struct {
		Server           string
		Group            string
		EventChangeTopic string
		EventChangeSet   string
	}
	Data toolkit.M
}

func NewAppConfig() *AppConfig {
	a := new(AppConfig)
	a.Hosts = make(map[string]string)
	a.Connections = make(map[string]struct {
		Txt      string
		PoolSize int
	})
	return a
}

func (cfg *AppConfig) DataToEnv() {
	for k, v := range cfg.Data {
		switch v.(type) {
		case string:
			os.Setenv(k, v.(string))
		}
	}
}

/*
func GetConfig(eventServer, topic string, serve bool, s *kaos.Service) (*AppConfig, error) {
	ev := knats.NewEventHub(eventServer, byter.NewByter(""))
	defer ev.Close()

	cfg := new(AppConfig)
	e := ev.Publish(topic, "", cfg)
	if e == nil && serve {
		go ServeConfigChange(cfg, s)
	}
	return cfg, e
}

func ServeConfigChange(cfg *AppConfig, s *kaos.Service) {
}

func ServeConfigSet(cmd, key string, cfg *AppConfig, s *kaos.Service) {
}
*/

func GetConfigFromEventHub(ev kaos.EventHub, topic string) (*AppConfig, error) {
	res := new(AppConfig)
	if e := ev.Publish(topic, "", res); e != nil {
		return nil, fmt.Errorf("fail get config from nats server. %s", e.Error())
	}
	return res, nil
}
