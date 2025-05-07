package dmutils

import (
	"net"
	"time"
)

func Telnet(nType, dns string) error {
	telnet, err := net.DialTimeout(nType, dns, 1*time.Second)
	if err != nil {
		return err
	}
	if telnet != nil {
		defer telnet.Close()
	}
	return nil
}
