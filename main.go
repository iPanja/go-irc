package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go-irc/format"
	"go-irc/parser"
	"net"
	"os"
	"strings"
	"time"
)

// const address = "irc.w3.org:6679"
// const address = "irc.libera.chat:6697"
const address = "irc.freenode.net:6667"
const nick = "tst45"

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func main() {
	i := make(chan string, 1)
	o := make(chan string)

	go handleUserInput(i)
	go handleServerResponse(i, o)
	go handleServerConnection(nick, i, o)

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
		sanitized := strings.Trim(raw, " \r\n")

		if len(sanitized) != 0 {
			i <- format.FormatUserInput(sanitized)
		}
	}
}

func handleServerResponse(i chan string, o chan string) {
	// No need for select, this thread can be blocking
	for {
		line := <-o

		// Respond to the periodic PING/PONG request
		if strings.HasPrefix(line, "PING") {
			// channel `i` is unbuffered
			// therefore when the send happens below, we are guaranteed it is received before execution resumes
			// aka this thread blocks until someone RECEIVES our item
			code := parser.ExtractPingPongCode(line)
			i <- fmt.Sprintf("PONG %s", code)

			fmt.Println(Yellow, "PING PONG!")
			continue
		}

		i := parser.NewIRCMessage(line)

		fmt.Println(i.FormatMessage())
	}
}

func handleServerConnection(nick string, i chan string, o chan string) {
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
		_, err = conn.Write([]byte(msg))
		if err != nil {
			panic(err)
		}
	}

	buf := make([]byte, 512) // An IRC message is at most 512 bytes; It will end with an `\r\n`
	tmp := make([]byte, 512)
	n := 0

MainLoop:
	for {
		// Check for user input (to send) without blocking reader
		select {
		case msg := <-i:
			fmt.Println(Green, "SENDING:", msg)
			conn.Write([]byte(fmt.Sprintf("%s\r\n", msg)))
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second))
			n, err = conn.Read(tmp)

			if err != nil && !os.IsTimeout(err) {
				break MainLoop
			}
		}

		if n == 0 {
			continue
		}

		buf = append(buf, tmp[:n]...)
		if bytes.HasSuffix(buf, []byte("\r\n")) || bytes.HasSuffix(buf, []byte("\n")) || len(buf) == 512 {
			o <- string(bytes.Trim(buf, "\x00\r\n"))
			buf = make([]byte, 512)
		}
	}

	// Error
	fmt.Println("Buf: ", buf)
	fmt.Println("Tmp: ", tmp)
	fmt.Println("Error: ", err)
	panic(err)
}
