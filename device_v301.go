package optocontrol

import (
	"sync"
	"time"

	"github.com/tarm/serial"
)

type device_v301 struct {
	conf *serial.Config
	port *serial.Port
	mux  sync.Mutex
}

type Device_v301 interface {
	Open() error
	Close() error
	ReadData() (*DataV3_01, error)
}

func NewV301(port string, speed int, timeout time.Duration) Device_v301 {
	conf := serial.Config{
		Name:        port,
		Baud:        speed,
		ReadTimeout: timeout,
	}

	dev := &device_v301{}
	dev.conf = &conf
	dev.mux = sync.Mutex{}

	return dev
}

func (dev *device_v301) Open() error {

	s, err := serial.OpenPort(dev.conf)
	if err != nil {
		return err
	}

	dev.port = s

	return nil
}

func (dev *device_v301) Close() error {

	return dev.port.Close()
}
