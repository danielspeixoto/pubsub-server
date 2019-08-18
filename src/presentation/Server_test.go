package presentation

import (
	"log"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestUDPServer(t *testing.T) {
	port := 10001
	go StartServer(port)
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
		msg := "sub"
		buf := []byte(msg)
		_, err = Conn.Write(buf)
		if err != nil {
			t.Fatal(err)
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
			response := string(buffer[:n])
			log.Println("msg received from server is: " + response)
			if "started" != response {
				t.Fatal("state not received properly")
			}
			time.Sleep(time.Second * 1)
			i++
		}
	})

	t.Run("Multiple clients can subscribe to state changes", func(t *testing.T) {
		clientCount := 3
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
				msg := "sub"
				buf := []byte(msg)
				_, err = Conn.Write(buf)
				if err != nil {
					t.Fatal(err)
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
					response := string(buffer[:n])
					log.Println("msg received from server is: " + response)
					if "started" != response {
						t.Fatal("state not received properly")
					}
					time.Sleep(time.Second * 1)
					j++
				}
			}()
			time.Sleep(time.Second * 3)
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

		msg := "state update:ok"
		buf := []byte(msg)
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
		msg = "sub"
		buf = []byte(msg)
		_, err = Conn.Write(buf)
		if err != nil {
			t.Fatal(err)
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
			response := string(buffer[:n])
			log.Println("msg received from server is: " + response)
			if "ok" != response {
				t.Fatal("state not received properly")
			}
			time.Sleep(time.Second * 1)
			i++
		}
	})
}
