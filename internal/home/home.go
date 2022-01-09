package home

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"

	"github.com/maytanthegeek/sync-rgb/pkg/tuya"
)

type HomeDevices struct {
	Strip TuyaDeviceLight
}

func NewHomeDevices(config string) *HomeDevices {
	content, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	payload := []map[string]string{}
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	dm := tuya.NewDeviceManagerRaw()
	// dm.DisableLogging()

	s := new(ITuyaDeviceLight)
	dm.DefineDevice(payload[0]["name"], payload[0]["gwId"], payload[0]["key"], payload[0]["ip"], payload[0]["version"], s)

	d, _ := dm.GetDevice("monitor")
	monitor := d.(TuyaDeviceLight)
	monitor.SetW("21", "colour", 5*time.Second)

	return &HomeDevices{Strip: monitor}
}
