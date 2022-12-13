package main

import (
	m "machine"
	"sync"

	"github.com/s-fairchild/reef-controller/comms"
	"github.com/s-fairchild/reef-controller/waterlevel"
)

var wg sync.WaitGroup

func main() {
	defer wg.Done()
	// Aquarium water level
	wl := waterlevel.NewWaterLevelSensor(m.GPIO17, m.GPIO15, m.GPIO14, m.PinInputPullup)
	wl.InitWaterLevel()
	wl.InitSignalLeds(m.GPIO13, m.GPIO12, m.GPIO11)

	wg.Add(1)
	go wl.MonitorLevel()
	wg.Wait()
	println("Water level monitor go routine exited")
}

func init() {
	err := comms.InitUART(m.UART0)
	if err != nil {
		println(err)
	}
}
