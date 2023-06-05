package kconfigurator

import (
	"time"

	"git.kanosolution.net/kano/kaos"
)

type configuratorModel struct {
	LastUpdate time.Time
	appConfig  *AppConfig
}

func NewConfigurator(cfg *AppConfig) *configuratorModel {
	m := new(configuratorModel)
	m.LastUpdate = time.Now()
	m.appConfig = cfg
	return m
}

func (c *configuratorModel) Read(ctx *kaos.Context, any string) (*AppConfig, error) {
	return c.appConfig, nil
}

func (c *configuratorModel) Write(ctx *kaos.Context, cfg *AppConfig) (*AppConfig, error) {
	c.LastUpdate = time.Now()
	c.appConfig = cfg
	ev, _ := ctx.DefaultEvent()
	ev.Publish(c.appConfig.EventServer.EventChangeTopic, c.appConfig, nil, nil)
	return c.appConfig, nil
}

type Request struct {
	Kind  string `json:"kind"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (c *configuratorModel) Set(ctx *kaos.Context, r *Request) (string, error) {
	c.LastUpdate = time.Now()
	switch r.Kind {
	case "Host":
		c.appConfig.Hosts[r.Key] = r.Value
	case "Hub":
		h, ok := c.appConfig.Connections[r.Key]
		if !ok {
			h = struct {
				Txt      string
				PoolSize int
			}{
				Txt: r.Value, PoolSize: 20,
			}
		}
		h.Txt = r.Value
		c.appConfig.Connections[r.Key] = h
	}

	ev, _ := ctx.DefaultEvent()
	ev.Publish(c.appConfig.EventServer.EventChangeTopic, c.appConfig, nil, nil)
	return "OK", nil
}

type configuratorEvent struct {
	appConfig *AppConfig
}

func (ce *configuratorEvent) Get(ctx *kaos.Context, parm string) (*AppConfig, error) {
	return ce.appConfig, nil
}

func (c *configuratorModel) EventModel() *configuratorEvent {
	em := new(configuratorEvent)
	em.appConfig = c.appConfig
	return em
}
