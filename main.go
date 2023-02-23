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
	// c.Rtc.SetTime(time.Date(2023, 02, 23, 15, 15, 00, 00, time.UTC))
	t, err := c.Rtc.ReadTime()
	if err != nil {
		panic(err)
	}
	println("Current Time:", t.Format(time.RFC3339))

	wl := waterlevel.New(m.GP17, m.GP15, m.LED, *c)
	wl.Init()
	go wl.MonitorLevel()

	magnesiumPump := dosing.New(m.GP18, dosing.Magnesium, c.Rtc)
	err = magnesiumPump.Configure(&dosing.DosingConfig{
		Ml:       30,
		Interval: 24 * time.Hour,
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
