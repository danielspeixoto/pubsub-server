package src

import (
	"bytes"
	"encoding/gob"
	"log"
	"net"
	"strconv"
	"time"
)

type Message struct {
	Timestamp int64
	Type      string
	Message   string
}

// Converts an interface to byte array
func GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Converts byte array to Message
func GetMessage(buff []byte) (Message, error) {
	dec := gob.NewDecoder(bytes.NewReader(buff))
	var msg Message
	err := dec.Decode(&msg)
	return msg, err
}

func StartServer(port int, period time.Duration) {
	packetConn, err := net.ListenPacket("udp", ":"+strconv.Itoa(port))
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
	// Periodically send status update to all subs
	go func() {
		for {
			log.Println("sending state to subs")
			buff, err := GetBytes(Message{
				Timestamp: time.Now().Unix(),
				Type:      "state",
				Message:   state,
			})
			for _, addr := range addrs {
				log.Println("sending package to " + addr.String())
				_, err = packetConn.WriteTo(buff, addr)
				if err != nil {
					log.Fatal(err)
				}
			}
			time.Sleep(period)
		}
	}()

	// Receives status updates from clients
	for {
		buffer := make([]byte, 1024)
		n, addr, err := packetConn.ReadFrom(buffer)
		if err != nil {
			log.Println(err)
			continue
		}
		msg, err := GetMessage(buffer[:n])
		if err != nil {
			log.Println(err)
			continue
		}
		// There are two types of message
		if msg.Type == "sub" {
			// Client subscribes to state changes
			if contains(addrs, addr) {
				log.Println("client: " + addr.String() + " is already  a sub")
			} else {
				addrs = append(addrs, addr)
			}
			log.Println("sending sub confirmation to " + addr.String())
		} else if msg.Type == "update" {
			// Client updates state
			state = msg.Message
		} else {
			// [Danger Area] Be Aware of possible attacks or versioning Problems
			log.Println("should not be here! Addr: " +
				addr.String() + ", bytes: " + string(buffer))
		}
		// Confirms message successfully delivered
		buff, err := GetBytes(Message{
			Timestamp: time.Now().Unix(),
			Type:      "ack",
			Message:   "ok",
		})
		_, err = packetConn.WriteTo(buff, addr)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}

func contains(s []net.Addr, e net.Addr) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
