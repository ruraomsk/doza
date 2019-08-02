package dozimetr

import (
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/mikepb/go-serial"
)

//Dozimetr одно измерение от дозиметра
type Dozimetr struct {
	Ntime   int     //Время накопления дозы пульта
	Value   float64 //Текущеее значение дозы пульта
	Pogr    float64 //Погрешность в % текущего значения
	SumDoza float64 //Накопленная доза за время Ntime
}

//GetTime return hour,minutes and secs
func (d *Dozimetr) GetTime() (hour int, min int, sec int) {
	hour = d.Ntime / 3600
	min = (d.Ntime - hour*3600) / 60
	sec = d.Ntime - (hour*3600 + min*60)
	return hour, min, sec
}
func getBit(xd byte, n uint) int {
	m := 1 << n
	if xd&byte(m) != 0 {
		return 1
	}
	return 0
}
func xD(xd byte) (izB byte, izC float64, izD string) {
	izB = byte((getBit(xd, 7) << 1) | (getBit(xd, 6)))
	c := (getBit(xd, 5) << 1) | (getBit(xd, 4))
	d := (getBit(xd, 3) << 3) | (getBit(xd, 2) << 2) | (getBit(xd, 1) << 1) | (getBit(xd, 0))
	switch d {
	case 0:
		izD = "%"
	case 1:
		izD = "%"
	case 2:
		izD = "мкЗв"
	case 3:
		izD = "мкЗв/ч"
	case 4:
		izD = "1(с*см^2)"
	default:
		izD = "********"
	}
	switch c {
	case 0:
		izC = 1.0
	case 1:
		izC = 1e3
	case 2:
		izC = 1e6
	default:
		izC = 1.0

	}
	return izB, izC, izD
}
func getFloat(buf []byte, pos int) float64 {
	result := 0.0
	s := 1
	if buf[pos] > 127 {
		s = -1
	}
	p := int16(buf[pos])
	if p > 127 {
		p -= 128
	}
	p -= 63

	m := (float64(uint16(buf[pos+1])*256+uint16(buf[pos+2])) / 65536) + 1
	result = float64(s) * m * math.Pow(2.0, float64(p))
	return result
}
func openPort(namePort string) (*serial.Port, error) {
	option := serial.RawOptions
	option.BitRate = 9600
	option.Parity = serial.PARITY_NONE
	option.DataBits = 8
	option.StopBits = 2
	option.RTS = serial.RTS_ON
	option.DTR = serial.DTR_OFF
	port, err := option.Open(namePort)
	return port, err
}
func oneByte(port *serial.Port) (byte, error) {
	b := make([]byte, 1)
	for true {
		nBytes, err := port.Read(b)
		if err != nil {
			if err.Error() == "EOF" {
				// fmt.Print("eof")
				continue
			}
			fmt.Println(err.Error(), port.Info.Name())
			return 0, err
		}
		if nBytes < 1 {
			continue
		}
		// fmt.Print(b[0], " ")
		return b[0], nil
	}
	return 0, nil
}
func exit(c chan *Dozimetr) {
	c <- nil
}
func crcCalc(buf []byte) bool {
	var crc uint16
	crc = 0xffff
	for i := 0; i < len(buf)-2; i++ {
		crcL := (crc & 0xff) ^ uint16(buf[i])
		crc = (crc & 0xff00) | (crcL & 0xff)
		for j := 0; j < 8; j++ {
			if crc&1 > 0 {
				crc ^= 0xa001
			}

		}
	}
	fmt.Println(crc&0xff, (crc&0xff00)>>8)
	if (crc&0xff != uint16(buf[len(buf)-2])) || ((crc&0xff00)>>8 != uint16(buf[len(buf)-1])) {
		return false
	}
	return true
}

//RoutDozimetr read and conver data from dozimetr
func RoutDozimetr(c chan *Dozimetr, namePort string) {
	defer exit(c)
	list, err := serial.ListPorts()
	if err != nil {
		fmt.Println("Serial error ", err.Error())
		return
	}
	if namePort == "" {
		if runtime.GOOS == "linux" {
			for _, p := range list {
				if p.USBProduct() == "USB-Serial Controller" {
					namePort = p.Name()
					break
				}
			}
		} else {
			namePort = "COM3"
		}
		if namePort == "" {
			fmt.Println("Dozimetr error not found USB-Serial Controller")
			return

		}
	}
	port, err := openPort(namePort)
	if err != nil {
		fmt.Println("Serial error ", err.Error(), namePort)
		return
	}
	defer port.Close()
	t := time.Now()
	t = t.Add(time.Duration(10 * time.Second))
	port.SetReadDeadline(t)
	buffer := make([]byte, 24)
	port.ResetInput()
	for true {
		b, err := oneByte(port)
		if err != nil {
			fmt.Println(err.Error(), namePort)
			port.Close()
			port, err = openPort(namePort)
			if err != nil {
				fmt.Println("Serial error ", err.Error(), namePort)
				return
			}
			port.ResetInput()
			continue
		}
		if b != 1 {
			// port.ResetInput()
			continue
		}
		buffer[0] = 1
		b, err = oneByte(port)
		if err != nil {
			fmt.Println(err)
			port.Close()
			port, err = openPort(namePort)
			if err != nil {
				fmt.Println("Serial error ", err.Error(), namePort)
				return
			}
			port.ResetInput()
			continue
		}
		if b != 12 {
			port.ResetInput()
			continue
		}
		buffer[1] = 12
		b, err = oneByte(port)
		if err != nil {
			// fmt.Println(err)
			port.Close()
			port, err = openPort(namePort)
			if err != nil {
				fmt.Println("Serial error ", err.Error(), namePort)
				return
			}
			port.ResetInput()
			continue
		}
		if b != 19 {
			port.ResetInput()
			continue
		}
		buffer[2] = 19
		buf := make([]byte, 21)
		nBytes, err := port.Read(buf)
		if err != nil {
			if err.Error() == "EOF" {
				// fmt.Println("eof")
				continue
			}
			fmt.Println(err.Error(), namePort)
			port.Close()
			port, err = openPort(namePort)
			if err != nil {
				fmt.Println("Serial error ", err.Error(), namePort)
				return
			}
			port.ResetInput()
			continue
		}
		// port.ResetInput()
		if nBytes < 1 {
			// fmt.Println("<0!")
			continue
		}
		// if nBytes != 21 {
		// 	port.ResetInput()
		// 	fmt.Println("not 21 length")
		// }
		// fmt.Println(buf)
		for i := 0; i < len(buf); i++ {
			buffer[3+i] = buf[i]
		}
		// if !crcCalc(buffer) {
		// fmt.Printf("port.Read: %v \n", buffer)
		// }

		// 6 байт это вроде P
		_, pC, _ := xD(buffer[6])
		// fmt.Println("6 byte ", pB, pC, pD)
		ntime := uint(buffer[18])*3600 + uint(buffer[19])*60 + uint(buffer[20])
		_, nC, _ := xD(buffer[6])
		d := new(Dozimetr)
		d.Ntime = int(ntime)
		d.Value = getFloat(buffer, 3) * pC
		d.SumDoza = getFloat(buffer, 14) * nC
		d.Pogr = getFloat(buffer, 7)
		c <- d
		// fmt.Println("send to chan")
		t := time.Now()
		t = t.Add(time.Duration(10 * time.Second))
		// port.Close()
		// port, err = openPort(namePort)
		port.SetReadDeadline(t)
	}
	c <- nil
}
