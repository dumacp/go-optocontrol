package optocontrol

import (
	"encoding/binary"
	"fmt"
)

type DataV3_01 struct {
	AdultUpFront   uint
	AdultDownFront uint
	KidUpFront     uint
	KidDownFront   uint
	LocksFront     uint
	AdultUpBack    uint
	AdultDownBack  uint
	KidUpBack      uint
	KidDownBack    uint
	LocksBack      uint

	RestartFront uint
	RestartBack  uint
	Raw          []byte
}

func (d *DataV3_01) RawResponse() []byte {
	return d.Raw
}

func (dev *device_v301) ReadData() (*DataV3_01, error) {

	dev.mux.Lock()
	defer dev.mux.Unlock()

	cmd := []byte{0x02, 0x51, 03}
	// cmd = append(cmd, parity(cmd))
	cmd = append(cmd, xor(cmd))

	fmt.Printf("request: [% X]\n", cmd)
	if _, err := dev.port.Write(cmd); err != nil {
		return nil, err
	}

	buf := make([]byte, 64)
	n, err := dev.port.Read(buf)
	if err != nil {
		return nil, err
	}

	if n > 1 && buf[1] == 'S' && n <= 32 {
		return nil, fmt.Errorf("data isn't complete, data: [% X], len: %d", buf[:n], len(buf[:n]))
	} else if n > 1 && buf[1] == 'N' {
		fmt.Printf("without changes: [% X]\n", buf[:n])
		return nil, nil
	}
	if buf[n-1] != 0x5A && xor(buf[:n-1]) != buf[n-1] {
		return nil, fmt.Errorf("parity error: %X != %X, [%X]", buf[n-1], xor(buf[:n-1]), buf[:n])
	}
	if buf[n-2] != 0x03 && buf[0] != 0x02 {
		return nil, fmt.Errorf("frame error: [%X], len: %d", buf[:n], len(buf[:n]))
	}

	data := &DataV3_01{
		AdultUpFront:   uint(binary.LittleEndian.Uint16(buf[4:6])),
		KidUpFront:     uint(binary.LittleEndian.Uint16(buf[6:8])),
		AdultDownFront: uint(binary.LittleEndian.Uint16(buf[8:10])),
		KidDownFront:   uint(binary.LittleEndian.Uint16(buf[10:12])),

		LocksFront:   uint(binary.LittleEndian.Uint16(buf[12:14])),
		RestartFront: uint(binary.LittleEndian.Uint16(buf[16:18])),

		AdultUpBack:   uint(binary.LittleEndian.Uint16(buf[18:20])),
		KidUpBack:     uint(binary.LittleEndian.Uint16(buf[20:22])),
		AdultDownBack: uint(binary.LittleEndian.Uint16(buf[22:24])),
		KidDownBack:   uint(binary.LittleEndian.Uint16(buf[24:26])),

		LocksBack:   uint(binary.LittleEndian.Uint16(buf[26:28])),
		RestartBack: uint(binary.LittleEndian.Uint16(buf[30:32])),

		// Locks:     uint(binary.LittleEndian.Uint16(buf[8:10])),
	}

	cpBuff := make([]byte, n)
	copy(cpBuff, buf)
	data.Raw = cpBuff

	return data, nil
}

func xor(data []byte) byte {
	sum := byte(0)
	for _, v := range data[:] {
		sum ^= byte(v)
	}
	return sum
}
