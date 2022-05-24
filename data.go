package optocontrol

import (
	"encoding/binary"
	"fmt"
)

type Data struct {
	AdultUp   uint
	AdultDown uint
	KidUp     uint
	KidDown   uint
	Locks     uint
}

type DataType int

func (d DataType) ADDR() byte {
	switch d {
	case DOOR_DATA_1:
		return 0x10
	case DOOR_DATA_2:
		return 0x20
	}
	return 0x00
}

const (
	DOOR_DATA_1 DataType = iota
	DOOR_DATA_2
)

func (dev *device) ReadData(addr DataType) (*Data, error) {

	dev.mux.Lock()
	defer dev.mux.Unlock()

	cmd := []byte{0x01, 0x52, addr.ADDR(), 0x0A, 0x00}
	// cmd = append(cmd, parity(cmd))
	cmd = append(cmd, 0x5A)

	fmt.Printf("request: [% X]\n", cmd)
	if _, err := dev.port.Write(cmd); err != nil {
		return nil, err
	}

	buf := make([]byte, 0x0A)
	n, err := dev.port.Read(buf)
	if err != nil {
		return nil, err
	}

	if n < 0x0A {
		return nil, fmt.Errorf("data isn't complete, data: [% X]", buf[:n])
	}
	if buf[n-1] != 0x5A && parity(buf[:n]) != buf[n-1] {
		return nil, fmt.Errorf("parity error: %X != %X", buf[n-1], parity(buf[:n]))
	}

	data := &Data{
		AdultUp:   uint(binary.LittleEndian.Uint16(buf[0:2])),
		AdultDown: uint(binary.LittleEndian.Uint16(buf[2:4])),
		KidUp:     uint(binary.LittleEndian.Uint16(buf[4:6])),
		KidDown:   uint(binary.LittleEndian.Uint16(buf[6:8])),
		Locks:     uint(binary.LittleEndian.Uint16(buf[8:10])),
	}

	return data, nil
}

func (dev *device) ReadSensor(addr SensorType) (*Sensors, error) {

	dev.mux.Lock()
	defer dev.mux.Unlock()

	cmd := []byte{0x01, 0x52, addr.ADDR(), 0x06}
	if _, err := dev.port.Write(cmd); err != nil {
		return nil, err
	}

	buf := make([]byte, 10)
	n, err := dev.port.Read(buf)
	if err != nil {
		return nil, err
	}

	if n < 6 {
		return nil, fmt.Errorf("data isn't complete, data: [% X]", buf[:n])
	}
	if buf[n-1] != 0x5A && parity(buf[:n]) != buf[n-1] {
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

func parity(data []byte) byte {
	sum := 0
	for _, v := range data[:len(data)-1] {
		sum += int(v)
	}
	return byte((sum % 0x100) & 0xFF)
}
