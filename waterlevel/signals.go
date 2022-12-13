package waterlevel

import (
	m "machine"
	"time"
)

func (w *waterLevel) InitSignalLeds(emptyReservoir, delay, noError m.Pin) {
	emptyReservoir.Configure(m.PinConfig{Mode: m.PinOutput})
	delay.Configure(m.PinConfig{Mode: m.PinOutput})
	noError.Configure(m.PinConfig{Mode: m.PinOutput})

	w.emptyReservoirLed = emptyReservoir
	w.delayLed = delay
	w.noError = noError
}

func (w *waterLevel) checkStatusAll() {
	if !w.delayLed.Get() && !w.emptyReservoirLed.Get() {
		w.noError.High()
	} else {
		w.noError.Low()
	}
}

func (w *waterLevel) emptySignal() {
	for w.reservoir.Get() {
		w.emptyReservoirLed.High()
		time.Sleep(500 * time.Millisecond)

		w.emptyReservoirLed.Low()
		time.Sleep(500 * time.Millisecond)
	}
}
