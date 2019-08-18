package presentation

import (
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func StartServer(port int) {
	packetConn, err := net.ListenPacket("udp", ":" + strconv.Itoa(port))
	log.Println("server running on port " + string(port))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := packetConn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	state := "started"
	addrs := make([]net.Addr, 0)
	log.Println("list of subs created and empty")

	go func() {
		for {
			log.Println("sending state to subs")
			for _, addr := range addrs {
				log.Println("sending package to " + addr.String())
				_, err := packetConn.WriteTo([]byte(state), addr)
				if err != nil {
					log.Fatal(err)
				}
			}
			time.Sleep(time.Second * 1)
		}
	}()

	for {
		buffer := make([]byte, 1024)
		n, addr, err := packetConn.ReadFrom(buffer)
		if err != nil {
			continue
		}
		msg := string(buffer[:n])
		// using this approach may result in concurrency issues
		// state and addrs is shared among goroutines
		if msg == "sub" {
			addrs = append(addrs, addr)
		} else if strings.HasPrefix(msg, "state update:") {
			state = msg[13:]
		}
	}
}