package waterlevel

import (
	"fmt"
	m "machine"
	"time"
)

type waterLevel struct {
	waterLevel  m.Pin // Pump sensor
	pumpRelay m.Pin // Pump relay
	reservoir m.Pin // Reservoir sensor
	sensorMode m.PinMode // should be PinInputPullup
	pumpDelay time.Time // actual time is meaningless unless an RTC is added to track time across power cycles
	// LED Signals
	emptyReservoirLed m.Pin
	delayLed m.Pin
	noError m.Pin
}

var (
	volumePumped float32 // volume pumped since last delay
	totalVolumePumped float32
)

// gallons per second
const (
	gps float32 = 0.017
)

type WaterLevel interface {
	InitWaterLevel()
	MonitorLevel()
	InitSignalLeds(emptyReservoir, delay, noError m.Pin)
}

func NewWaterLevelSensor(pumpSensorPin, pumpRelayPin, reservoirPin m.Pin, mode m.PinMode) WaterLevel {
	return &waterLevel{
		waterLevel:  pumpSensorPin,
		pumpRelay: pumpRelayPin,
		reservoir: reservoirPin,
		sensorMode: mode,
	}
}

// InitWaterLevel configures the water level sensor pin and relay pin
func (w *waterLevel) InitWaterLevel() {
	fmt.Printf("Initializing water level sensor on pin %d in mode %d\n", w.waterLevel, w.sensorMode)
	fmt.Printf("Pump flow rate is %f gallons per second\n", gps)
	w.waterLevel.Configure(m.PinConfig{Mode: w.sensorMode})
	w.pumpRelay.Configure(m.PinConfig{Mode: m.PinOutput})
	w.reservoir.Configure(m.PinConfig{Mode: w.sensorMode})
	m.LED.Configure(m.PinConfig{Mode: m.PinOutput})
	w.pumpDelay = time.Now()
}

// MonitorLevel polls the sensor pin status once per second to determine if the water level has dropped below the sensor.
//
// If the sensor returns false, the water pump is actuated.
//
// A maximum of 1 gallon per 12 hours is set to prevent overflowing.
func (w *waterLevel) MonitorLevel() {
	println("Starting water level sensor monitoring")
	m.LED.High()
	for {
		go w.checkStatusAll() // TODO troubleshoot why this doesn't work outside of the loop
		if !w.waterLevel.Get() {
			w.checkReservoirLevel()
			w.actuatePumpRelay()
		} else {
			w.pumpRelay.Low()
		}
		time.Sleep(1 * time.Second)
	}
}

// actuatePumpRelay checks if the total volume pumped is greater than 1 gallon.
//
// If more than one gallon has been pumped, a 12 hour delay is set.
//
// After 24 hours the pump can be activated again.
//
// Without an RTC the time delay is lost across power cycles
func (w *waterLevel) actuatePumpRelay() {
	if volumePumped >= 1.0 {
		fmt.Printf("%f water pumped, shutting off pump for time delay\n", volumePumped)
		w.pumpRelay.Low()
		fmt.Printf("Pump time delay: %v\n", w.pumpDelay)
		w.pumpDelay = time.Now().Add(12 * time.Hour)
		// capture total volume to display on something like an oled screen
		totalVolumePumped = volumePumped
		volumePumped = 0.0
	}

	if w.pumpDelay.Before(time.Now()) {
		w.pumpRelay.High()
		w.delayLed.Low()
		volumePumped += gps
		fmt.Printf("Water pump is on\nGallons pumped: %f\n", volumePumped)
	} else {
		w.delayLed.High()
	}
}

// checkReservoirLevel blocks until the reservoir isn't empty
func (w *waterLevel) checkReservoirLevel() {
	if w.reservoir.Get() {
		w.pumpRelay.Low()
		fmt.Printf("Reservoir empty\n")
		w.emptySignal()
	} else {
		w.emptyReservoirLed.Low()
	}
}
