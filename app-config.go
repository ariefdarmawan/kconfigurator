package kconfigurator

import (
	"fmt"
	"os"
	"strings"

	"git.kanosolution.net/kano/kaos"
	"github.com/sebarcode/codekit"
)

var (
	Variables codekit.M
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

func (cfg *AppConfig) Parse() {
	for id, host := range cfg.Hosts {
		cfg.Hosts[id] = Update(host)
	}

	for id, data := range cfg.Data {
		dataStr, ok := data.(string)
		if ok {
			cfg.Data[id] = Update(dataStr)
		}
	}

	for id, conn := range cfg.Connections {
		conn.Txt = Update(conn.Txt)
		cfg.Connections[id] = conn
	}

	cfg.EventServer.Server = Update(cfg.EventServer.Server)
	cfg.EventServer.Group = Update(cfg.EventServer.Group)
	cfg.EventServer.EventChangeSet = Update(cfg.EventServer.EventChangeSet)
	cfg.EventServer.EventChangeTopic = Update(cfg.EventServer.EventChangeTopic)
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
	if e := ev.Publish(topic, "", res, nil); e != nil {
		return nil, fmt.Errorf("fail get config from nats server. %s", e.Error())
	}
	return res, nil
}

func UpdateWithEnv(txt string) string {
	parts := strings.Split(txt, "${env:")
	if len(parts) <= 1 {
		return txt
	}

	for _, part := range parts[1:] {
		envID := strings.Split(part, "}")[0]
		envValue := os.Getenv(envID)
		txt = strings.ReplaceAll(txt, "${env:"+envID+"}", envValue)
	}

	return txt
}

func UpdateWithVar(txt string) string {
	parts := strings.Split(txt, "${var:")
	if len(parts) <= 1 {
		return txt
	}

	for _, part := range parts[1:] {
		id := strings.Split(part, "}")[0]
		value, ok := Variables[id].(string)
		if ok {
			txt = strings.ReplaceAll(txt, "${var:"+id+"}", value)
		}
	}

	return txt
}

func Update(txt string) string {
	s := UpdateWithEnv(txt)
	if s == "" || s == txt {
		s = UpdateWithVar(txt)
	}

	// special command
	wd, e := os.Getwd()
	if e == nil && strings.Contains(s, "${ctx:wd}") {
		s = strings.ReplaceAll(s, "${ctx:wd}", wd)
	}

	return s
}
