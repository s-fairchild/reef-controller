package main

import (
	m "machine"
	"time"

	"github.com/s-fairchild/reef-controller/comms"
	"github.com/s-fairchild/reef-controller/dosing"
	"github.com/s-fairchild/reef-controller/rtc"
	"github.com/s-fairchild/reef-controller/waterlevel"
)

func main() {
	c := rtc.New(&m.I2CConfig{
		SDA: m.GP26,
		SCL: m.GP27,
	}, m.I2C1)
	c.Init()
	c.Rtc.Write([]byte(time.Now().Format(rtc.LayoutTime)))
	t, err := c.Rtc.ReadTime()
	if err != nil {
		panic(err)
	}
	println("Current Time:", t.Format(rtc.LayoutDate))

	wl := waterlevel.New(m.GP17, m.GP15, m.LED, *c)
	wl.Init()
	go wl.MonitorLevel()

	magnesiumPump := dosing.New(m.GP18, "magnesium-pump", c.Rtc)
	err = magnesiumPump.Configure(&dosing.DosingConfig{
		Ml:       30,
		Interval: 24,
	})
	if err != nil {
		panic(err)
	}

	go magnesiumPump.Dose()

	select {}
}

func init() {
	err := comms.InitUART(m.UART0, true)
	if err != nil {
		println(err)
	}
}
