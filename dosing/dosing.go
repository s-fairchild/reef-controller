package dosing

import (
	"errors"
	m "machine"
	"time"

	"github.com/s-fairchild/reef-controller/rtc"

	"tinygo.org/x/drivers/ds1307"
)

type dosingPump struct {
	name         string
	pump         m.Pin
	config       *DosingConfig
	sram         ds1307.Device
	offset       int64
	bytesWritten int // TODO possibly remove this if unused
	lastRun      time.Time
}

type DosingPump interface {
	Configure(c *DosingConfig) error
	Dose() error
}

// DosingConfig holds the dosing pump configuration
//
// ml is the volume of liquid to be dispend
//
// Hour, Minute, and Second are the time of day to dose in 24 hour format
// Interval is how often to dose within a 24 hour period
type DosingConfig struct {
	Ml uint8
	// Hour	 int
	// Minute   int
	// Second   int
	Interval time.Duration
}

func New(pump m.Pin, name string, sram ds1307.Device) DosingPump {
	return &dosingPump{
		pump: pump,
		name: name,
		sram: sram,
	}
}

func (d *dosingPump) Configure(c *DosingConfig) error {
	if c.Ml == 0 {
		return errors.New("dosing pump config Ml cannot be 0")
	}

	d.config = c
	d.pump.Configure(m.PinConfig{Mode: m.PinOutput})
	d.pump.Low()
	var err error
	d.offset, err = d.sram.Seek(0, 0)
	if err != nil {
		return err
	}

	println("Dosing pump will dose", c.Ml, "ml's every", d.config.Interval, "hours")
	return nil
}

// Dose
func (d *dosingPump) Dose() error {
	println("Starting dosing pump", d.name)
	tick := time.NewTicker(time.Duration(d.config.Ml) * time.Second)
	var err error
	for {
		println("entering loop")
		r := make([]uint8, 40, 40)
		// _, err = rtc.ReadSavedTime(r, d.offset, d.sram)
		_, err = rtc.ReadSavedTime(r, d.offset, d.sram)
		if err != nil {
			return err
		}
		println("past readtime")

		lastRun, err := time.Parse(rtc.LayoutTime, string(r))
		if err != nil {
			return err
		}
		println(d.name, "last run", string(r))

		if lastRun.After(lastRun) {
			println("Activating dosing pump", d.name, "now")
			d.pump.High()
			<-tick.C
			d.pump.Low()
			println("Deactivating dosing pump", d.name, "now")

			d.bytesWritten, err = rtc.WriteTime(time.Now(), d.bytesWritten, d.offset, d.sram)
			if err != nil {
				return err
			}

		}
	}
}
