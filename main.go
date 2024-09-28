package main

import (
	"bufio"
	"fmt"
	"go-irc/format"
	"net"
	"os"
	"strings"
	"time"
)

// const address = "irc.w3.org:6679"
// const address = "irc.libera.chat:6697"
const address = "irc.freenode.net:6667"
const nick = "tst45"

func main() {
	var i chan string = make(chan string)

	go handleUserInput(i)
	go handleServerConnection(nick, i)

	for {
		// block forever until CTRL+C
		<-make(chan bool)
	}
}

func handleUserInput(i chan string) {
	r := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		raw, _ := r.ReadString('\n')
		sanitized := strings.ToUpper(strings.Trim(raw, " \n\r"))

		if len(sanitized) != 0 {
			i <- format.FormatUserInput(sanitized)
		}
	}
}

func handleServerConnection(nick string, i chan string) {
	d := net.Dialer{Timeout: 5 * time.Second}
	conn, err := d.Dial("tcp", address)
	defer conn.Close()
	if err != nil {
		panic(err)
	}

	// Initiate handshake
	messages := []string{
		fmt.Sprintf("NICK %s\r\n", nick),
		fmt.Sprintf("USER %s 0 * :IRCTesting\r\n", nick),
	}
	for _, msg := range messages {
		_, err := conn.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}

	buf := make([]byte, 512)
	tmp := make([]byte, 256)
	n := 0

	for {
		// Check for user input (to send) without blocking reader
		select {
		case msg := <-i:
			conn.Write([]byte(fmt.Sprintf("%s\r\n", msg)))
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, err = conn.Read(tmp)

			if err != nil && !os.IsTimeout(err) {
				fmt.Println("Buf: ", buf)
				fmt.Println("Tmp: ", tmp)
				fmt.Println("Error: ", err)
				panic(err)
			}
		}

		response := string(tmp[:n])
		// Respond to the periodic PING/PONG request
		if strings.HasPrefix(response, "PING") {
			code := strings.Split(response, " ")[1]
			conn.Write([]byte(fmt.Sprintf("PONG %s", code)))
			fmt.Println("PING PONG!")
			continue
		}

		fmt.Print(response)
		buf = append(buf, tmp[:n]...)
	}
}
