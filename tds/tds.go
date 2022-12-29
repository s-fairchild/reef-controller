package tds

import (
	"errors"
	"machine"
	"time"
)

type device struct {
	adc            machine.ADC
	averageVoltage float32
	aref           float32 // ADC Reference voltage in millivolts
	resolution     float32
}

type Device interface {
	// Read collects readings from ADC and converts the analog voltage read into Total Dissolved Solids ppm.
	// Temperature provided in celsius is used to calculate the temperature compensation.
	//
	// Note: 1.8 is the closest for this manufacture. If no temperature can be provided, passing a constant 65.0 will result in 1.8 being calculated.
	Read(temp float32) (float32, error)
}

const (
	readCycle         time.Duration = time.Millisecond * 40
	sampleCount                     = 30
	defaultReference                = 5.0 // 5.0v
	defaultResolution               = 1023.0
)

// New returns a new total dissolve solids sensor driver given an ADC pin.
func New(adc machine.ADC, aref, resolution float32) Device {
	return &device{
		adc:        adc,
		aref:       aref,
		resolution: resolution,
	}
}

func (t *device) Read(temp float32) (float32, error) {
	rs := t.collectSamples()
	m := t.findMedian(rs)
	tds := t.voltage2tds(t.adc2voltage(m), temp)
	if tds < 0.0 || tds > 1000.0 {
		return tds, errors.New("total dissolved solids reading is invalid:")
	}
	return tds, nil
}

func (t *device) collectSamples() []uint16 {
	readBuffer := make([]uint16, sampleCount)
	for i := 0; i < len(readBuffer); i++ {
		readBuffer[i] = t.adc.Get()
		time.Sleep(readCycle)
	}
	return readBuffer
}

func (t *device) findMedian(d []uint16) uint16 {
	t.bubbleSort(d)

	var median uint16
	l := len(d)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		median = (d[l/2-1] + d[l/2]) / 2
	} else {
		median = d[l/2]
	}
	return median
}

func (d *device) adc2voltage(n uint16) float32 {
	return (float32(n) * d.aref) / d.resolution
}

// voltage2tds converts averaged voltage value to tds value
// formulas were converted from source wiki Arduino example: http://www.cqrobot.wiki/index.php/TDS_(Total_Dissolved_Solids)_Meter_Sensor_SKU:_CQRSENTDS01#Arduino_Application
func (t *device) voltage2tds(v float32, temp float32) float32 {
	compVolt := t.calcVoltCompensation(v, t.calcTempCompCoefficient(temp))
	return (133.42*compVolt*compVolt*compVolt - 255.86*compVolt*compVolt + 857.39*compVolt) * 0.5
}

// calcTempCompCoefficient calculates the compensation for temperature differences
// Should be very close to 1.8
func (t *device) calcTempCompCoefficient(temp float32) float32 {
	return 1.0 + 0.02*(temp-25.0)
}

func (t *device) calcVoltCompensation(v float32, tempCompCo float32) float32 {
	return v / tempCompCo
}

func (t *device) bubbleSort(d []uint16) {
	for i := 0; i < len(d)-1; i++ {
		for j := 0; j < len(d)-i-1; j++ {
			if d[j] > d[j+1] {
				d[j+1], d[j] = d[j], d[j+1]
				break
			}
		}
	}
}
