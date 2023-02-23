package dosing

import (
	"errors"
	m "machine"
	"time"

	"github.com/s-fairchild/reef-controller/rtc"

	"tinygo.org/x/drivers/ds1307"
)

type dosingPump struct {
	name         Liquid
	pump         m.Pin
	config       *DosingConfig
	sram         ds1307.Device
	offset       int64
	bytesWritten uint8 // used to make slice size when reading time
	lastRun      time.Time
}

type DosingPump interface {
	Configure(c *DosingConfig) error
	Dose() error
}

type Liquid string

const (
	Magnesium Liquid = "magnesium"
)

type DosingConfig struct {
	// ml is the volume of liquid to be dispend
	Ml       uint8
	// Hour, Minute, and Second are the time of day to dose in 24 hour format
	// Interval is how often to dose within a 24 hour period
	Interval time.Duration // TODO Make interval configurable with button and LCD user interface
}

func New(pump m.Pin, liquid Liquid, sram ds1307.Device) DosingPump {
	return &dosingPump{
		pump: pump,
		name: liquid,
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

	s, err := rtc.ReadSavedTime(35, 0, d.sram)
	if err != nil {
		return err
	}
	println("====================================================================")
	println("===================== Configuration ================================")
	println("====================================================================")
	println("saved time read to save to SRAM:", s.LastDose.Format(time.RFC3339))
	println("Dosing pump will dose", d.config.Ml, "ml's every", int(d.config.Interval.Hours()), "hours")
	println("====================================================================")

	// Set last dose if one doesn't exist.
	// This will result in LastDose being the same as current time, causing the Dosing loop
	// to wait until Interval is over to actually dose, typically 24 hours
	if s.LastDose.IsZero() {
		println(err)
		t, err := d.sram.ReadTime()
		if err != nil {
			return err
		}

		println("Last dose not found, setting now as last dose time.")
		d.bytesWritten, err = rtc.WriteTime(t, d.offset, d.sram)
		if err != nil {
			println("failed to initial time")
			panic(err)
		}
	}

	return nil
}

// Dose
func (d *dosingPump) Dose() error {
	println("Starting dosing pump", d.name)
	for {
		s, err := rtc.ReadSavedTime(d.bytesWritten, d.offset, d.sram)
		if err != nil {
			return err
		}

		t, err := d.sram.ReadTime()
		if err != nil {
			return err
		}
		println("===================== Dosing - ", d.name, " ========================")
		println("====================================================================")
		println("last run", s.LastDose.Format(time.RFC3339))
		println("Current time:", t.Format(time.RFC3339))
		println("====================================================================")

		if t.After(s.LastDose) {
			d.pump.High()
			println("Activating dosing pump", d.name, "now")

			time.Sleep(time.Duration(d.config.Ml) * time.Second)

			d.pump.Low()
			println("Deactivating dosing pump", d.name, "now")

			t, err := d.sram.ReadTime()
			if err != nil {
				return err
			}

			d.bytesWritten, err = rtc.WriteTime(t, d.offset, d.sram)
			if err != nil {
				return err
			}
		}
		println("Waiting", int(d.config.Interval.Hours()), "hours before dosing...")
		time.Sleep(d.config.Interval)
	}
}
