package src

import (
	"log"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestUDPServer(t *testing.T) {
	port := 10001
	go StartServer(port, time.Second)
	t.Run("A client can subscribe to state changes", func(t *testing.T) {
		ServerAddr, err :=
			net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(port))
		if err != nil {
			t.Fatal(err)
		}

		LocalAddr, err :=
			net.ResolveUDPAddr("udp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}

		Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
		if err != nil {
			t.Fatal(err)
		}
		defer Conn.Close()

		i := 0
		msg := Message{
			Timestamp: time.Now().Unix(),
			Type:      "sub",
			Message:   []byte(""),
		}
		buf, err := GetBytes(msg)
		if err != nil {
			t.Fatal(err)
		}
		_, err = Conn.Write(buf)
		if err != nil {
			t.Fatal(err)
		}
		buffer := make([]byte, 1024)
		n, _, err := Conn.ReadFromUDP(buffer)
		if err != nil {
			t.Fatal(err)
		}
		response, err := GetMessage(buffer[:n])
		if err != nil {
			t.Fatal(err)
		}
		log.Println("msg received from server is: " + string(response.Message))
		if "ok" != string(response.Message) {
			t.Fatal("ack not received properly")
		}
		for {
			if i > 3 {
				break
			}
			buffer := make([]byte, 1024)
			n, _, err := Conn.ReadFromUDP(buffer)
			if err != nil {
				t.Fatal(err)
			}
			response, err := GetMessage(buffer[:n])
			if err != nil {
				t.Fatal(err)
			}
			log.Println("msg received from server is: " + string(response.Message))
			if "started" != string(response.Message) {
				t.Fatal("state not received properly")
			}
			time.Sleep(time.Second * 1)
			i++
		}
	})

	t.Run("Multiple clients can subscribe to state changes", func(t *testing.T) {
		clientCount := 20000
		i := 0
		for {
			if i > clientCount {
				break
			}
			go func() {
				log.Println("client number " + strconv.Itoa(i))
				ServerAddr, err :=
					net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(port))
				if err != nil {
					t.Fatal(err)
				}

				LocalAddr, err :=
					net.ResolveUDPAddr("udp", "127.0.0.1:0")
				if err != nil {
					t.Fatal(err)
				}

				Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
				if err != nil {
					t.Fatal(err)
				}
				defer Conn.Close()

				j := 0
				msg := Message{
					Timestamp: time.Now().Unix(),
					Type:      "sub",
					Message:   []byte(""),
				}
				buf, err := GetBytes(msg)
				if err != nil {
					t.Fatal(err)
				}
				_, err = Conn.Write(buf)
				if err != nil {
					t.Fatal(err)
				}
				buffer := make([]byte, 1024)
				n, _, err := Conn.ReadFromUDP(buffer)
				if err != nil {
					t.Fatal(err)
				}
				response, err := GetMessage(buffer[:n])
				if err != nil {
					t.Fatal(err)
				}
				log.Println("msg received from server is: " + string(response.Message))
				if "ok" != string(response.Message) {
					t.Fatal("ack not received properly")
				}
				for {
					if j > 3 {
						break
					}
					buffer := make([]byte, 1024)
					n, _, err := Conn.ReadFromUDP(buffer)
					if err != nil {
						// May Fail because it is UDP ¯\_(ツ)_/¯
						t.Fatal(err)
					}
					response, err := GetMessage(buffer[:n])
					if err != nil {
						t.Fatal(err)
					}
					log.Println("msg received from server is: " + string(response.Message))
					if "started" != string(response.Message) {
						t.Fatal("state not received properly")
					}
					time.Sleep(time.Second * 1)
					j++
				}
			}()
			time.Sleep(time.Millisecond * 50)
			i++
		}
	})

	t.Run("A client can update the server state", func(t *testing.T) {
		ServerAddr, err :=
			net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(port))
		if err != nil {
			t.Fatal(err)
		}

		LocalAddr, err :=
			net.ResolveUDPAddr("udp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}

		Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
		if err != nil {
			t.Fatal(err)
		}
		defer Conn.Close()

		msg := Message{
			Timestamp: time.Now().Unix(),
			Type:      "update",
			Message:   []byte("ok"),
		}
		buf, err := GetBytes(msg)
		if err != nil {
			t.Fatal(err)
		}
		_, err = Conn.Write(buf)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 1)

		// new subscription
		ServerAddr, err =
			net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(port))
		if err != nil {
			t.Fatal(err)
		}

		LocalAddr, err =
			net.ResolveUDPAddr("udp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}

		Conn, err = net.DialUDP("udp", LocalAddr, ServerAddr)
		if err != nil {
			t.Fatal(err)
		}
		defer Conn.Close()

		i := 0
		msg = Message{
			Timestamp: time.Now().Unix(),
			Type:      "sub",
			Message:   []byte(""),
		}
		buf, err = GetBytes(msg)
		if err != nil {
			t.Fatal(err)
		}
		_, err = Conn.Write(buf)
		if err != nil {
			t.Fatal(err)
		}
		buffer := make([]byte, 1024)
		n, _, err := Conn.ReadFromUDP(buffer)
		if err != nil {
			t.Fatal(err)
		}
		response, err := GetMessage(buffer[:n])
		if err != nil {
			t.Fatal(err)
		}
		log.Println("msg received from server is: " + string(response.Message))
		if "ok" != string(response.Message) {
			t.Fatal("ack not received properly")
		}
		for {
			if i > 3 {
				break
			}
			buffer := make([]byte, 1024)
			n, _, err := Conn.ReadFromUDP(buffer)
			if err != nil {
				t.Fatal(err)
			}
			response, err := GetMessage(buffer[:n])
			if err != nil {
				t.Fatal(err)
			}
			log.Println("msg received from server is: " + string(response.Message))
			if "ok" != string(response.Message) {
				t.Fatal("state not received properly")
			}
			time.Sleep(time.Second * 1)
			i++
		}
	})

	t.Run("server does not stop when one of the subs fails to connect", func(t *testing.T) {
		go func() {
			ServerAddr, err :=
				net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(port))
			if err != nil {
				t.Fatal(err)
			}

			LocalAddr, err :=
				net.ResolveUDPAddr("udp", "127.0.0.1:0")
			if err != nil {
				t.Fatal(err)
			}

			Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
			if err != nil {
				t.Fatal(err)
			}
			msg := Message{
				Timestamp: time.Now().Unix(),
				Type:      "sub",
				Message:   []byte(""),
			}
			buf, err := GetBytes(msg)
			if err != nil {
				t.Fatal(err)
			}
			_, err = Conn.Write(buf)
			if err != nil {
				t.Fatal(err)
			}
			buffer := make([]byte, 1024)
			n, _, err := Conn.ReadFromUDP(buffer)
			if err != nil {
				t.Fatal(err)
			}
			response, err := GetMessage(buffer[:n])
			if err != nil {
				t.Fatal(err)
			}
			log.Println("msg received from server is: " + string(response.Message))
			if "ok" != string(response.Message) {
				t.Fatal("ack not received properly")
			}
			err = Conn.Close()
			if err != nil {
				t.Fatal(err)
			}
			log.Println("closed")
		}()
		time.Sleep(time.Second * 1)
		ServerAddr, err :=
			net.ResolveUDPAddr("udp", "127.0.0.1:"+strconv.Itoa(port))
		if err != nil {
			t.Fatal(err)
		}

		LocalAddr, err :=
			net.ResolveUDPAddr("udp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}

		Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
		if err != nil {
			t.Fatal(err)
		}
		defer Conn.Close()

		i := 0
		msg := Message{
			Timestamp: time.Now().Unix(),
			Type:      "sub",
			Message:   []byte(""),
		}
		buf, err := GetBytes(msg)
		if err != nil {
			t.Fatal(err)
		}
		_, err = Conn.Write(buf)
		if err != nil {
			t.Fatal(err)
		}
		buffer := make([]byte, 1024)
		n, _, err := Conn.ReadFromUDP(buffer)
		if err != nil {
			t.Fatal(err)
		}
		response, err := GetMessage(buffer[:n])
		if err != nil {
			t.Fatal(err)
		}
		log.Println("msg received from server is: " + string(response.Message))
		if "ok" != string(response.Message) {
			t.Fatal("ack not received properly")
		}
		for {
			if i > 3 {
				break
			}
			buffer := make([]byte, 1024)
			n, _, err := Conn.ReadFromUDP(buffer)
			if err != nil {
				t.Fatal(err)
			}
			response, err := GetMessage(buffer[:n])
			if err != nil {
				t.Fatal(err)
			}
			log.Println("msg received from server is: " + string(response.Message))
			if "started" != string(response.Message) {
				t.Fatal("state not received properly")
			}
			time.Sleep(time.Second * 1)
			i++
		}
	})
}
