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
	wl := waterlevel.NewWaterLevelSensor(m.GPIO17, m.GPIO15, m.PinInputPullup)
	wl.InitWaterLevel()

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
