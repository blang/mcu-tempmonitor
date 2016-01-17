package main

import (
	"net"
	"strconv"
)

type HandlerFunc func(float32)

func ListenAndServeUDP(listen string, handler HandlerFunc) error {
	addr, err := net.ResolveUDPAddr("udp", ":10001")
	if err != nil {
		return err
	}

	ServerConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer ServerConn.Close()

	buf := make([]byte, 64)
	for {
		n, _, err := ServerConn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		i, err := strconv.Atoi(string(buf[0:n]))
		if err != nil {
			continue
		}
		handler(float32(i) / 10000)
	}
	return nil
}
