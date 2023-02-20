package main

import (
	m "machine"

	"github.com/s-fairchild/reef-controller/comms"
	"github.com/s-fairchild/reef-controller/rtc"
	"github.com/s-fairchild/reef-controller/waterlevel"
)

func main() {
	c := rtc.New(&m.I2CConfig{
		SDA: m.GP26,
		SCL: m.GP27,
	}, m.I2C1)
	c.Init()
	t, err := c.Rtc.ReadTime()
	if err != nil {
		panic(err)
	}
	println("Current Time:", t.Format(rtc.Format))

	wl := waterlevel.NewWaterLevelSensor(m.GP17, m.GP15, m.LED, c)

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
