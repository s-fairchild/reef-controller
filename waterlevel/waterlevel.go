package waterlevel

import (
	m "machine"
	"time"

	"github.com/s-fairchild/reef-controller/rtc"
)

type waterLevel struct {
	waterLevel m.Pin     // Pump sensor
	pumpRelay  m.Pin     // Pump relay
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

type WaterLevel interface {
	// InitWaterLevel configures the water level sensor pin and relay pin
	Init()
	// MonitorLevel polls the sensor pin status once per second to determine if the water level has dropped below the sensor.
	//
	// If the sensor returns false, the water pump is actuated.
	//
	// A maximum of 1 gallon per 12 hours is set to prevent overflowing.
	MonitorLevel()
}

func New(pumpSensorPin, pumpRelayPin m.Pin, led m.Pin, rtc rtc.Rtc) WaterLevel {
	return &waterLevel{
		waterLevel: pumpSensorPin,
		pumpRelay:  pumpRelayPin,
		mLED:       led,
		clock:      rtc,
	}
}

func (w *waterLevel) Init() {
	println("Initializing water level sensor on pin ", w.waterLevel)
	println("Pump flow rate is ", gps, " gallons per second")
	w.waterLevel.Configure(m.PinConfig{Mode: m.PinInputPullup})
	w.pumpRelay.Configure(m.PinConfig{Mode: m.PinOutput})
	w.pumpDelay = time.Now()
	w.mLED.Configure(m.PinConfig{Mode: m.PinOutput})
	w.mLED.High()
}

func (w *waterLevel) MonitorLevel() {
	println("Starting water level sensor monitoring")
	for {
		// TODO Make pumping stop when water level sensor is true again
		// Currently once it starts pumping, it won't stop until (previously the reservoir was empty) 1 gallon has been pumped
		if !w.waterLevel.Get() {
			w.actuatePumpRelay()
		} else {
			w.pumpRelay.Low()
		}
		time.Sleep(1 * time.Second)
	}
}

// TODO use channels to start threads for actuating the pump. If the reservoir goes empty while the pump is on, kill the routine.
// actuatePumpRelay checks if the total volume pumped is greater than 1 gallon.
//
// If more than one gallon has been pumped, a 12 hour delay is set.
//
// After 24 hours the pump can be activated again.
func (w *waterLevel) actuatePumpRelay() error {
	if volumePumped >= 1.0 {
		println(volumePumped, " water pumped, shutting off pump for time delay")
		w.pumpRelay.Low()
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
		w.pumpRelay.High()
		volumePumped += gps
		println("Water pump is on\nGallons pumped:", volumePumped)
	} else {
	}

	return nil
}
