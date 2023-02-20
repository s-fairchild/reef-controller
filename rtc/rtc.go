package rtc

import (
	"errors"
	m "machine"

	"tinygo.org/x/drivers/ds1307"
)

// format is used to pass to time.Format
const Format = "2006-01-02 15:04:05"

type Rtc struct {
	config *m.I2CConfig
	i2c    *m.I2C
	Rtc    ds1307.Device
}

func New(config *m.I2CConfig, i2c *m.I2C) *Rtc {
	return &Rtc{
		config: config,
		i2c:    i2c,
	}
}

func (c *Rtc) Init() error {
	c.i2c.Configure(*c.config)
	c.Rtc = ds1307.New(c.i2c)
	if !c.Rtc.IsOscillatorRunning() {
		return errors.New("failed to initialize rtc, oscillator is not running")
	}
	return nil
}
