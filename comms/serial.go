package comms

import (
	"errors"
	m "machine"
)

// InitUART configures the provided UART device and sets it to the default UART when def is true.
//
// Set default as false for UART's intended for data only communications, rather than text console output
func InitUART(u *m.UART, def bool) error {
	println("Initializing UART output")
	err := u.Configure(m.UARTConfig{})
	if err != nil {
		return errors.New("failed to initialize UART0:" + err.Error())
	}
	println("Successfully initialized UART serial output")

	// println and fmt will use the default UART for output
	if u != m.DefaultUART && def {
		m.DefaultUART = u
	}

	return nil
}
