package main

import (
	m "machine"

	"github.com/s-fairchild/reef-controller/comms"
	"github.com/s-fairchild/reef-controller/waterlevel"
)

func main() {
	// Aquarium water level
	wl := waterlevel.NewWaterLevelSensor(m.GP17, m.GP15, m.LED)
	wl.Init()

	go wl.MonitorLevel()

	select {}
}

func init() {
	err := comms.InitUART(m.UART0, true)
	if err != nil {
		println(err)
	}
}
