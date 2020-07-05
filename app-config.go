package kconfigurator

import (
	"git.kanosolution.net/kano/kaos"
	"git.kanosolution.net/kano/kaos/events/knats"
	"github.com/ariefdarmawan/byter"
	"github.com/eaciit/toolkit"
)

type AppConfig struct {
	Hosts       map[string]string
	Connections map[string]struct {
		Txt      string
		PoolSize int
	}
	EventServers struct {
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
