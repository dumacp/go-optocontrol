package optocontrol

import (
	"sync"
	"time"

	"github.com/tarm/serial"
)

type device struct {
	conf *serial.Config
	port *serial.Port
	mux  sync.Mutex
}

type Device interface {
	Open() error
	Close() error
	ReadData(addr DataType) (*Data, error)
	ReadSensors(addr SensorType) (*Sensors, error)
}

func New(port string, speed int, timeout time.Duration) Device {
	conf := serial.Config{
		Name:        port,
		Baud:        speed,
		ReadTimeout: timeout,
	}

	dev := &device{}
	dev.conf = &conf
	dev.mux = sync.Mutex{}

	return dev
}

func (dev *device) Open() error {

	s, err := serial.OpenPort(dev.conf)
	if err != nil {
		return err
	}

	dev.port = s

	return nil
}

func (dev *device) Close() error {

	return dev.port.Close()
}
