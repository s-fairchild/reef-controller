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
	bytesWritten uint8 // used to make slice size when reading time
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

	s, err := rtc.ReadSavedTime(32, 0, d.sram)
	if err != nil {
		return err
	}
	println("====================================================================")
	println("===================== Configuration ================================")
	println("====================================================================")
	println("saved time read to save to SRAM:", s.LastDose.Format(time.RFC3339))
	println("Dosing pump will dose", c.Ml, "ml's every", d.config.Interval, "hours")
	println("====================================================================")

	// See if a previous dose has ran before powering on
	// The current time must be used in place of the last dose if not
	_, err = time.Parse(time.RFC3339, s.LastDose.Format(time.RFC3339))
	if err != nil {
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
	tick := time.NewTicker(time.Duration(d.config.Ml) * time.Second)
	tickDay := time.NewTicker(24 * time.Hour)
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
			println("Activating dosing pump", d.name, "now")
			d.pump.High()
			<-tick.C
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
		<-tickDay.C
	}
}
