package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"
)

type Proceso struct {
	Id        int
	Count     int
	Is_closed bool
}

var process_slide []*Proceso

func printSlide() {
	for {
		fmt.Println("--------------------")
		for _, p := range process_slide {
			fmt.Println("id ", p.Id, ": ", p.Count)
			p.Count += 1
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func servidor() {
	go printSlide()
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleClient(c)

	}
}

func handleClient(c net.Conn) {
	var p Proceso
	err := gob.NewDecoder(c).Decode(&p)
	if err != nil {
		fmt.Println(err)
		return
	} else if p.Is_closed == false {
		c2, err2 := net.Dial("tcp", ":9998")
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		err3 := gob.NewEncoder(c2).Encode(process_slide[0])
		if err3 != nil {
			fmt.Println(err3)
		}
		c2.Close()
		process_slide = append(process_slide[:0], process_slide[1:]...)
	} else {
		process_slide = append(process_slide, &p)
	}
}


func main() {
	for i := 1; i <= 6; i++ {
		process_slide = append(process_slide, &Proceso{
			Id:        i,
			Count:     0,
			Is_closed: false,
		})
	}
	go servidor()

	var input string
	fmt.Scanln(&input)
}
