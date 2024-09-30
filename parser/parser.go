package parser

import (
	"strings"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Magenta = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

type Parser interface {
	ExtractData(buf []byte) map[string]string
}

type IRCMessage struct {
	Prefix  string   // The sender (user or server)
	Command string   // (PRIVMSG, PING, JOIN, etc)
	Params  []string // (channel, message content, etc)

	Raw string
}

func NewIRCMessage(raw string) IRCMessage {
	result := IRCMessage{Raw: raw}

	parts := strings.Split(raw, "\x20") // space

	if strings.HasPrefix(raw, "\x3a") { // colon
		result.Prefix = parts[0]
		parts = parts[1:]
	}

	result.Command = parts[0]
	result.Params = parts[1:]

	return result
}

func ExtractPingPongCode(s string) string {
	parts := strings.Split(strings.Trim(s, "\x00\r\n"), " ")
	code := strings.Split(parts[1], "\n")[0]

	return code
}

func (m IRCMessage) FormatMessage() string {
	color := White

	switch m.Command {
	case "PRIVMSG":
		color = Magenta
	case "372", "005", "251", "253", "254", "255", "265", "266", "375":
		color = Gray
	case "NOTICE":
		color = Cyan
	default:
		color = White
	}

	return color + " " + m.Raw
}
