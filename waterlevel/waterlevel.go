package waterlevel

import (
	"fmt"
	m "machine"
	"time"
)

type waterLevel struct {
	sensorPin  m.Pin
	relayPin m.Pin
	sensorMode m.PinMode // should be PinInputPullup
}

type WaterLevel interface {
	InitWaterLevel()
	MonitorLevel()
}

func NewWaterLevelSensor(sPin, rPin m.Pin, mode m.PinMode) WaterLevel {
	return &waterLevel{
		sensorPin:  sPin,
		relayPin: rPin,
		sensorMode: mode,
	}
}

func (w *waterLevel) InitWaterLevel() {
	fmt.Printf("Initializing water level sensor on pin %d in mode %d\n", w.sensorPin, w.sensorMode)
	w.sensorPin.Configure(m.PinConfig{Mode: w.sensorMode})
	w.relayPin.Configure(m.PinConfig{Mode: m.PinOutput})
}

func (w *waterLevel) MonitorLevel() {
	println("Starting water level sensor monitoring")
	for {
		if !w.sensorPin.Get() {
			println("Water pump is on")
			w.relayPin.High()
		} else {
			w.relayPin.Low()
		}
		time.Sleep(1 * time.Second)
	}
}
