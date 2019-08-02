package main

import (
	"fmt"
	"net"
	"rura/doza/dozimetr"
	"rura/teprol/logger"
	"sync"
	"time"
)

type dozaValue struct {
	value string
	flag  string
}

var mutex sync.Mutex
var timer chan string

func sleeping() {
	for true {
		time.Sleep(3 * time.Second)
		timer <- "work"
	}

}
func connectToCombo() net.Conn {
	for true {
		conn, err := net.Dial("tcp", "192.168.10.30:5507")
		if err != nil {
			logger.Error.Println("Error tcp ", err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		logger.Info.Println("Connect to Combo")
		return conn
	}
	return nil
}
func main() {
	values := [2]dozaValue{{"0.0", "false"}, {"0.0", "false"}}
	names := [2]string{"/dev/ttyS1", "/dev/ttyS0"}
	Doza1 := make(chan *dozimetr.Dozimetr)
	Doza2 := make(chan *dozimetr.Dozimetr)
	timer = make(chan string)
	var conn net.Conn
	var err error
	err = logger.Init("/home/rura/log/doza")
	if err != nil {
		fmt.Println("Error init loggin subsystem ", err.Error())
		return
	}
	conn = connectToCombo()

	go sleeping()
	go dozimetr.RoutDozimetr(Doza1, names[0])
	go dozimetr.RoutDozimetr(Doza2, names[1])
	for true {
		select {
		case d := <-Doza1:
			{
				if d == nil {
					logger.Error.Println("End work " + names[0])
					mutex.Lock()
					values[0] = dozaValue{"0.0", "false"}
					mutex.Unlock()
					time.Sleep(1 * time.Second)
					go dozimetr.RoutDozimetr(Doza1, names[0])
					continue
				}
				// h, m, s := d.GetTime()
				st := fmt.Sprintf("%f", d.Value)
				mutex.Lock()
				if d.Value > 0.00001 {
					values[0] = dozaValue{st, "true"}
				}
				mutex.Unlock()
				logger.Info.Println(names[0], d.Value, d.Pogr, d.SumDoza)
				// isdata <- "data"
			}
		case d := <-Doza2:
			{
				if d == nil {
					logger.Error.Println("End work " + names[1])
					mutex.Lock()
					values[1] = dozaValue{"0.0", "false"}
					mutex.Unlock()
					time.Sleep(1 * time.Second)
					go dozimetr.RoutDozimetr(Doza2, names[1])
					continue
				}
				// h, m, s := d.GetTime()
				st := fmt.Sprintf("%f", d.Value)
				mutex.Lock()
				if d.Value > 0.00001 {
					values[1] = dozaValue{st, "true"}
				}
				mutex.Unlock()
				logger.Info.Println(names[1], d.Value, d.Pogr, d.SumDoza)
				// isdata <- "data"

			}
		case d := <-timer:
			{
				mutex.Lock()
				st := fmt.Sprint("[ ", values[0].value, " ", values[0].flag, " ", values[1].value, " ", values[1].flag, " ]\x00")
				if d == "work" {
					values[0] = dozaValue{"0.0", "false"}
					values[1] = dozaValue{"0.0", "false"}
				}
				mutex.Unlock()
				logger.Info.Println(d, st)
				for true {
					_, err = conn.Write([]byte(st))
					if err != nil {
						logger.Error.Println(err.Error())
						conn.Close()
						conn = connectToCombo()
					} else {
						break
					}

				}
			}

		}
	}
}
