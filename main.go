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
	wl := waterlevel.NewWaterLevelSensor(m.GP17, m.GP15, m.GP14, m.PinInputPullup)
	wl.InitWaterLevel()
	wl.InitSignalLeds(m.GP12, m.GP11, m.GP10)

	wg.Add(1)
	go wl.MonitorLevel()
	wg.Wait()
	println("Water level monitor go routine exited")
}

func init() {
	err := comms.InitUART(m.UART0, true)
	if err != nil {
		println(err)
	}
}
