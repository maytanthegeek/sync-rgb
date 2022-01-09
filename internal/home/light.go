package home

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/maytanthegeek/sync-rgb/pkg/tuya"
)

type TuyaDeviceLight interface {
	Set(string, interface{}) error
	SetW(string, interface{}, time.Duration) (interface{}, error)
	Status(string) interface{}
}

type ITuyaDeviceLight struct {
	tuya.BaseDevice
	// 20: power: true / false
	// 21: mode: white,colour,scene,music
	// 22: bright: 10-10000
	// 23: color temp: 10-10000
	// 24: color: h:0-360,s:0-1000,v:0-1000
	// 25: scene: 030e0d0000000000000001f401f4
	// 26: left time: 0
	// 28:
	status map[string]interface{}
}

func (s *ITuyaDeviceLight) Set(dps string, value interface{}) error {
	m := s.App.MakeBaseMsg()
	m["dps"] = map[string]interface{}{dps: value}
	return s.App.SendEncryptedCommand(tuya.CodeMsgSet, m)
}

func (s *ITuyaDeviceLight) SetW(dps string, value interface{}, delay time.Duration) (interface{}, error) {
	c := tuya.MakeSyncChannel()
	k := s.Subscribe(c)
	defer s.Unsubscribe(k)

	deadLine := time.Now().Add(delay * time.Second)
	m := s.App.MakeBaseMsg()
	m["dps"] = map[string]interface{}{dps: value}
	err := s.App.SendEncryptedCommand(tuya.CodeMsgSet, m)
	if err != nil {
		return nil, err
	}
	for {
		select {
		case synMsg := <-c:
			if synMsg.Code == tuya.CodeMsgStatus || synMsg.Code == tuya.CodeMsgAutoStatus {
				s := s.Status(dps)
				return s, err
			}
		case <-time.After(deadLine.Sub(time.Now())):
			return nil, errors.New("Timeout")
		}
	}
}

func (s *ITuyaDeviceLight) Status(key string) interface{} {
	if val, ok := s.status[key]; ok {
		return val
	} else {
		return nil
	}
}

func (s *ITuyaDeviceLight) ProcessResponse(code int, data []byte) {
	switch {
	case len(data) == 0:
		return
	case code == 7:
		return
	case code == 9:
		return
	}

	var r map[string]interface{}
	err := json.Unmarshal(data, &r)
	if err != nil {
		log.Println("JSON decode error")
		return
	}

	if val, ok := r["dps"]; ok {
		v2 := val.(map[string]interface{})
		for k, v := range v2 {
			s.status[k] = v
		}
	}
	s.Notify(code, s)
}

func (s *ITuyaDeviceLight) Configure(a *tuya.Appliance, c *tuya.ConfigurationData) {
	s.status = make(map[string]interface{})
	s.Init("TuyaDeviceLight", a, c)
}
