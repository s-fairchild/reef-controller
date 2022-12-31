package main

import (
	m "machine"

	"github.com/s-fairchild/reef-controller/comms"
	"github.com/s-fairchild/reef-controller/tds"
	"github.com/s-fairchild/reef-controller/waterlevel"
)

func main() {
	// Aquarium water level
	wl := waterlevel.NewWaterLevelSensor(m.GP17, m.GP15, m.GP14, m.PinInputPullup)
	wl.InitWaterLevel()
	wl.InitSignalLeds(m.GP26, m.GP27, m.GP28)

	s := tds.New(m.ADC0, 3.3, 65535.0)
	s.Configure()

	go wl.MonitorLevel()
	go func() {
		for {
			s.GetTds(20.0)
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
