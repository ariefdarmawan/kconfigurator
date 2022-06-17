package kconfigurator

import (
	"fmt"
	"os"

	"git.kanosolution.net/kano/kaos"
	"github.com/sebarcode/codekit"
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
	Data codekit.M
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

func GetConfigFromEventHub(ev kaos.EventHub, topic string) (*AppConfig, error) {
	res := new(AppConfig)
	if e := ev.Publish(topic, "", res); e != nil {
		return nil, fmt.Errorf("fail get config from nats server. %s", e.Error())
	}
	return res, nil
}
