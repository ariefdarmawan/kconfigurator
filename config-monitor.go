package kconfigurator

import (
	"errors"
	"fmt"

	"git.kanosolution.net/kano/dbflex"
	"git.kanosolution.net/kano/kaos"
	"github.com/ariefdarmawan/datahub"
	"github.com/eaciit/toolkit"
)

type configMonitor struct {
	cfg   *AppConfig
	svc   *kaos.Service
	topic string

	fnChange func()
}

func NewConfigMonitor(cfg *AppConfig, svc *kaos.Service) *configMonitor {
	cm := new(configMonitor)
	cm.cfg = cfg
	cm.svc = svc
	cm.topic = cfg.EventServer.EventChangeTopic
	return cm
}

/*
func (cm *configMonitor) OnConfigChanged(ev kaos.EventHub, s *kaos.Service) error {
	return ev.Subscribe(cm.topic, nil,
		func(ctx *kaos.Context, newCfg *AppConfig) (string, error) {
			cm.cfg = newCfg
			for connID := range newCfg.Connections {
				if h, e := MakeHub(newCfg, connID); e == nil {
					cm.svc.RegisterDataHub(h, connID)
				}
			}
			s.Log().Infof("Configuration has been changed")
			return "", nil
		})
}
*/

func MakeDbConn(config *AppConfig, s *kaos.Service) error {
	for k := range config.Connections {
		h, e := MakeHub(config, k)
		if e != nil {
			return fmt.Errorf("unable to initialize database %s: %s", k, e.Error())
		}
		s.RegisterDataHub(h, k)
	}
	return nil
}

func CloseDbConn(config *AppConfig, s *kaos.Service) {
	for k := range config.Connections {
		if h, e := s.GetDataHub(k); e == nil {
			s.Log().Infof("closing database connection %s", k)
			h.Close()
		}
	}
}

func MakeHub(config *AppConfig, name string) (*datahub.Hub, error) {
	ci, ok := config.Connections[name]
	if !ok {
		return nil, errors.New("Invalid ConnectionSettings " + name)
	}

	// intitate the hub
	h := datahub.NewHub(func() (dbflex.IConnection, error) {
		conn, err := dbflex.NewConnectionFromURI(ci.Txt, nil)
		if err != nil {
			return nil, err
		}
		if err = conn.Connect(); err != nil {
			return nil, fmt.Errorf("fail to connect. %s", err.Error())
		}
		conn.SetKeyNameTag("key")
		conn.SetFieldNameTag(toolkit.TagName())
		return conn, nil
	}, ci.PoolSize > 0, ci.PoolSize)

	// validate conn
	i, c, e := h.GetConnection()
	if e != nil {
		return nil, fmt.Errorf("fail to get connection. %s. ConnString: %s", e.Error(), ci.Txt)
	}
	h.CloseConnection(i, c)

	return h, nil
}
