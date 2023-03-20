package ato

import (
	m "machine"
	"time"

	"github.com/s-fairchild/reef-controller/rtc"
)

type ato struct {
	waterLevel m.Pin     // Pump sensor
	pump  m.Pin     // Pump relay
	pumpDelay  time.Time // actual time is meaningless unless an RTC is added to track time across power cycles
	mLED       m.Pin     // Machine LED
	clock      rtc.Rtc
}

var (
	volumePumped      float32 // volume pumped since last delay
	totalVolumePumped float32
)

const (
	// gallons per second
	gps float32 = 0.017
)

type Ato interface {
	// InitWaterLevel configures the water level sensor pin and relay pin
	Init()
	// MonitorLevel polls the sensor pin status once per second to determine if the water level has dropped below the sensor.
	//
	// If the sensor returns false, the water pump is actuated.
	//
	// A maximum of 1 gallon per 12 hours is set to prevent overflowing.
	MonitorLevel()
}

func New(pumpSensorPin, pumpRelayPin m.Pin, led m.Pin, rtc rtc.Rtc) Ato {
	return &ato{
		waterLevel: pumpSensorPin,
		pump:  pumpRelayPin,
		mLED:       led,
		clock:      rtc,
	}
}

func (w *ato) Init() {
	println("Initializing water level sensor on pin ", w.waterLevel)
	println("Pump flow rate is ", gps, " gallons per second")
	w.waterLevel.Configure(m.PinConfig{Mode: m.PinInputPullup})
	w.pump.Configure(m.PinConfig{Mode: m.PinOutput})
	w.pumpDelay = time.Now()
	w.mLED.Configure(m.PinConfig{Mode: m.PinOutput})
	w.mLED.High()
}

func (w *ato) MonitorLevel() {
	println("Starting water level sensor monitoring")
	for {
		if !w.waterLevel.Get() {
			w.actuatePump()
		} else {
			w.pump.Low()
		}
		time.Sleep(1 * time.Second)
	}
}

// TODO use channels to start threads for actuating the pump. If the reservoir goes empty while the pump is on, kill the routine.
// actuatePump checks if the total volume pumped is greater than 1 gallon.
//
// If more than one gallon has been pumped, a 12 hour delay is set.
//
// After 24 hours the pump can be activated again.
func (w *ato) actuatePump() error {
	if volumePumped >= 1.0 {
		println(volumePumped, " water pumped, shutting off pump for time delay")
		w.pump.Low()
		now, err := w.clock.Rtc.ReadTime()
		if err != nil {
			return err
		}
		w.pumpDelay = now.Add(12 * time.Hour)
		totalVolumePumped = volumePumped
		volumePumped = 0.0
	}

	now, err := w.clock.Rtc.ReadTime()
	if err != nil {
		return err
	}

	if w.pumpDelay.Before(now) {
		w.pump.High()
		volumePumped += gps
		println("Water pump is on\nGallons pumped:", volumePumped)
	} else {
	}

	return nil
}
