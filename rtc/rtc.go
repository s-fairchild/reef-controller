package rtc

import (
	"errors"
	m "machine"
	"time"

	"tinygo.org/x/drivers/ds1307"
)

const (
	LayoutTime   = "15:04:05"
	LayoutDate   = "2006-01-02 15:04:05"
	SRAMCapacity = 56 // SRAM capacity in bytes
)

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

// WriteTime encodes a time as a base64 string
// then writes the encoded time to sram
func WriteTime(t time.Time, offset int64, rtc ds1307.Device) (int, error) {
	rtc.Seek(offset, 0)

	println("time about to write to SRAM: ", t.Format(LayoutTime))
	println("Length of string written:", len([]byte(t.Format(LayoutTime))))
	b, err := rtc.Write([]byte(t.Format(LayoutTime)))
	if err != nil {
		return 0, err
	}
	if b != len([]byte(t.Format(LayoutTime))) {
		panic("bytes encoded not equal to bytes written to SRAM, failed to right dosing time to SRAM")
	}
	println("bytes written:", b)
	println("Wrote to SRAM:", t.Format(LayoutTime))
	return b, nil
}

func ReadSavedTime(data []byte, offset int64, rtc ds1307.Device) (int, error) {
	rtc.Seek(int64(offset), 0)

	b, err := rtc.Read(data)
	if err != nil {
		return 0, err
	}

	println("bytes read:", b)
	println("Read from SRAM:", string(data))

	if b != len(data) {
		return 0, errors.New("failed sanity check, Time bytes read from SRAM don't match bytes written")
	}
	return b, nil
}
