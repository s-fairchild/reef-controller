package comms

import (
	"fmt"
	m "machine"
)

func InitUART(u *m.UART) error {
	println("Initializing UART output")
	err := u.Configure(m.UARTConfig{})
	if err != nil {
		return fmt.Errorf("failed to initialize UART0, %s\n", err.Error())
	}
	println("Successfully initialized UART serial output")
	return nil
}
