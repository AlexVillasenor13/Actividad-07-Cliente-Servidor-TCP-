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

var p Proceso

func getCliente() {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(p)
	if err != nil {
		fmt.Println(err)
	}
	c.Close()
	go printProceso()
}

func printProceso() {
	for {
		if p.Is_closed {
			return
		}
		if p.Id != -1 {
			fmt.Println("id ", p.Id, ": ", p.Count)
			p.Count += 1
			time.Sleep(time.Millisecond * 500)
		}
	}
}

func listenCliente() {
	s, err := net.Listen("tcp", ":9998")
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
		err2 := gob.NewDecoder(c).Decode(&p)
		if err2 != nil {
			fmt.Println(err2)
			return
		}
		if p.Id != -1 {
			s.Close()
			return
		}
	}
}
func endCliente() {
	p.Is_closed = true
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = gob.NewEncoder(c).Encode(p)
	if err != nil {
		fmt.Println(err)
	}
	c.Close()
}

func main() {
	p.Id = -1

	go listenCliente()
	go getCliente()
	defer endCliente()
	var input string
	fmt.Scanln(&input)
}
