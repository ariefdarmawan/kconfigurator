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

func (cm *configMonitor) OnConfigChanged(ev kaos.EventHub, s *kaos.Service) error {
	return ev.Subscribe(cm.topic, s, nil,
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

func MakeHub(config *AppConfig, name string) (*datahub.Hub, error) {
	ci, ok := config.Connections[name]
	if !ok {
		return nil, errors.New("Invalid ConnectionSettings " + name)
	}
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
	return h, nil
}
