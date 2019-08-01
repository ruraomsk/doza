package main

import (
	"fmt"
	"rura/doza/dozimetr"
)

func main() {
	names := [7]string{"/dev/ttyS0", "/dev/ttyS1", "/dev/ttyS2", "/dev/ttyS3", "/dev/ttyS4", "/dev/ttyUSB0", "/dev/ttyUSB2"}
	Doza1 := make(chan *dozimetr.Dozimetr)
	Doza2 := make(chan *dozimetr.Dozimetr)
	Doza3 := make(chan *dozimetr.Dozimetr)
	Doza4 := make(chan *dozimetr.Dozimetr)
	Doza5 := make(chan *dozimetr.Dozimetr)
	Doza6 := make(chan *dozimetr.Dozimetr)
	Doza7 := make(chan *dozimetr.Dozimetr)
	go dozimetr.RoutDozimetr(Doza1, names[0])
	go dozimetr.RoutDozimetr(Doza2, names[1])
	go dozimetr.RoutDozimetr(Doza3, names[2])
	go dozimetr.RoutDozimetr(Doza4, names[3])
	go dozimetr.RoutDozimetr(Doza5, names[4])
	go dozimetr.RoutDozimetr(Doza6, names[5])
	go dozimetr.RoutDozimetr(Doza7, names[6])
	worked := 7
	for true {
		if worked == 0 {
			return
		}
		select {
		case d := <-Doza1:
			{
				if d == nil {
					fmt.Println("End work " + names[0])
					worked--
					continue
				}
				h, m, s := d.GetTime()
				fmt.Println(names[0], h, m, s, d.Value, d.Pogr, d.SumDoza)

			}
		case d := <-Doza2:
			{
				if d == nil {
					fmt.Println("End work " + names[1])
					worked--
					continue
				}
				h, m, s := d.GetTime()
				fmt.Println(names[1], h, m, s, d.Value, d.Pogr, d.SumDoza)

			}
		case d := <-Doza3:
			{
				if d == nil {
					fmt.Println("End work " + names[2])
					worked--
					continue
				}
				h, m, s := d.GetTime()
				fmt.Println(names[2], h, m, s, d.Value, d.Pogr, d.SumDoza)
			}
		case d := <-Doza4:
			{
				if d == nil {
					fmt.Println("End work " + names[3])
					worked--
					continue
				}
				h, m, s := d.GetTime()
				fmt.Println(names[3], h, m, s, d.Value, d.Pogr, d.SumDoza)
			}
		case d := <-Doza5:
			{
				if d == nil {
					fmt.Println("End work " + names[4])
					worked--
					continue
				}
				h, m, s := d.GetTime()
				fmt.Println(names[4], h, m, s, d.Value, d.Pogr, d.SumDoza)
			}
		case d := <-Doza6:
			{
				if d == nil {
					fmt.Println("End work " + names[5])
					worked--
					continue
				}
				h, m, s := d.GetTime()
				fmt.Println(names[5], h, m, s, d.Value, d.Pogr, d.SumDoza)
			}
		case d := <-Doza7:
			{
				if d == nil {
					fmt.Println("End work " + names[6])
					worked--
					continue
				}
				h, m, s := d.GetTime()
				fmt.Println(names[6], h, m, s, d.Value, d.Pogr, d.SumDoza)
			}

		}
	}
}
