package optocontrol

import "fmt"

type Sensors struct {
	Sensor1 int
	Sensor2 int
	Sensor3 int
	Sensor4 int
	Sensor5 int
	Sensor6 int
}

type SensorType int

const (
	DOOR_SENSOR_1 SensorType = iota
	DOOR_SENSOR_2
)

func (d SensorType) ADDR() byte {
	switch d {
	case DOOR_SENSOR_1:
		return 0x1a
	case DOOR_SENSOR_2:
		return 0x2a
	}
	return 0x00
}

const (
	DOOR_1 DataType = iota
	DOOR_2
)

func (dev *device) ReadSensors(addr SensorType) (*Sensors, error) {

	dev.mux.Lock()
	defer dev.mux.Unlock()

	cmd := []byte{0x01, 0x52, addr.ADDR(), 0x06}
	if _, err := dev.port.Write(cmd); err != nil {
		return nil, err
	}

	buf := make([]byte, 6)
	n, err := dev.port.Read(buf)
	if err != nil {
		return nil, err
	}

	if n < 6 {
		return nil, fmt.Errorf("data isn't complete, data: [% X]", buf[:n])
	}
	if parity(buf[:n]) != buf[n-1] {
		return nil, fmt.Errorf("parity error: %X != %X", buf[n-1], parity(buf[:n]))
	}

	data := &Sensors{
		Sensor1: int(buf[0]),
		Sensor2: int(buf[1]),
		Sensor3: int(buf[2]),
		Sensor4: int(buf[3]),
		Sensor5: int(buf[4]),
		Sensor6: int(buf[5]),
	}

	return data, nil
}
