package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// var CRLF = string([]byte{0x0D, 0x0A})
var CRLF = "\r\n"

// const address = "irc.w3.org:6679"
// const address = "irc.libera.chat:6697"
const address = "irc.freenode.net:6667"

func main() {
	var i chan string = make(chan string)

	go openInputReader(i)
	go openConnection(i)

	for {
		// block forever until CTRL+C
		<-make(chan bool)
	}
}

func preprocessUserInput(input string) string {
	if len(input) == 0 {
		return ""
	}

	if input[0] == '/' {
		parts := strings.Split(input[1:], " ")

		parts[0] = strings.ToUpper(parts[0])
		result := strings.Join(parts, " ")

		input = result
	}

	return input
}

func openInputReader(i chan string) {
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		raw, _ := r.ReadString('\n')
		sanitized := strings.ToUpper(strings.Trim(raw, " \n\r"))

		i <- preprocessUserInput(sanitized)
	}
}

func openConnection(i chan string) {
	d := net.Dialer{Timeout: 5 * time.Second}
	conn, err := d.Dial("tcp", address)
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("Established connection to IRC server")

	// Initiate handshake
	messages := []string{
		"NICK tst45\r\n",
		"USER tst45 0 * :IRCTesting\r\n",
	}
	for _, msg := range messages {
		_, err := conn.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}

	buf := make([]byte, 512)
	tmp := make([]byte, 256)
	for {
		select {
		case msg := <-i:
			conn.Write([]byte(fmt.Sprintf("%s\r\n", msg)))
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, err := conn.Read(tmp)

			if err != nil && !os.IsTimeout(err) {
				fmt.Println("Buf: ", buf)
				fmt.Println("Tmp: ", tmp)
				fmt.Println("Error: ", err)
				panic(err)
			}

			if n > 0 {
				response := string(tmp[:n])
				if strings.HasPrefix(response, "PING") {
					code := strings.Split(response, " ")[1]
					conn.Write([]byte(fmt.Sprintf("PONG %s", code)))
					fmt.Println("PING PONG!")
				} else {
					fmt.Print(response)
				}

				buf = append(buf, tmp[:n]...)
				tmp = make([]byte, 256)
			}
		}

	}
}
