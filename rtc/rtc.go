package rtc

import (
	"errors"
	m "machine"
	"time"

	"github.com/mailru/easyjson"
	"github.com/s-fairchild/reef-controller/types"
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
func WriteTime(t time.Time, offset int64, rtc ds1307.Device) (uint8, error) {
	rtc.Seek(offset, 0)

	d := &types.DosingState{
		LastDose: t,
	}
	bytes, err := d.MarshalJSON()
	if err != nil {
		return 0, err
	}

	println("time about to write to SRAM: ", t.Format(LayoutDate))
	b, err := rtc.Write(bytes)
	if err != nil {
		return 0, err
	}
	if b != len(bytes) {
		panic("bytes encoded not equal to bytes written to SRAM, failed to right dosing time to SRAM")
	}
	println("Wrote to SRAM:", t.Format(LayoutTime))
	return uint8(b), nil
}

func ReadSavedTime(size uint8, offset int64, rtc ds1307.Device) (*types.DosingState, error) {
	rtc.Seek(offset, 0)

	data := make([]uint8, size, size)
	b, err := rtc.Read(data)
	if err != nil {
		return &types.DosingState{}, err
	}

	state := &types.DosingState{}
	easyjson.Unmarshal(data, state)

	// println("Read from SRAM:", state.LastDose.Format(LayoutDate))

	if b != len(data) {
		return &types.DosingState{}, errors.New("failed sanity check, Time bytes read from SRAM don't match bytes written")
	}
	return state, nil
}
