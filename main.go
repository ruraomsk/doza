package main

import (
	"fmt"
	"rura/doza/dozimetr"
)

func main() {
	cDoza := make(chan *dozimetr.Dozimetr)
	go dozimetr.RoutDozimetr(cDoza)
	for true {
		d := <-cDoza
		if d == nil {
			fmt.Println("End work Dozimetr")
			return
		}
		h, m, s := d.GetTime()
		fmt.Println(h, m, s, d.Value, d.Pogr, d.SumDoza)
	}
}
