package main

import (
	m "machine"

	"github.com/s-fairchild/reef-controller/comms"
	"github.com/s-fairchild/reef-controller/tds"
	"github.com/s-fairchild/reef-controller/waterlevel"
)

const normTemp float32 = 25.5556 // 78F

func main() {
	// Aquarium water level
	wl := waterlevel.NewWaterLevelSensor(m.GP17, m.GP15, m.GP14, m.PinInputPullup)
	wl.InitWaterLevel()
	wl.InitSignalLeds(m.GP26, m.GP27, m.GP28)

	adc := m.ADC{Pin: m.ADC0}
	adc.Configure(m.ADCConfig{})
	s := tds.New(adc, 3.3, 65535.0)

	go wl.MonitorLevel()
	go func() {
		for {
			s.Read(normTemp)
		}
	}()

	select {}
}

func init() {
	err := comms.InitUART(m.UART0, true)
	if err != nil {
		println(err)
	}
}
