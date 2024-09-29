package parser

import "strings"

type Parser interface {
	ExtractData(buf []byte) map[string]string
}

type IRCMessage struct {
	Prefix  string   // The sender (user or server)
	Command string   // (PRIVMSG, PING, JOIN, etc)
	Params  []string // (channel, message content, etc)

	Code int
	Raw  string
}

func NewIRCMessage(raw string) IRCMessage {
	if len(raw) == 0 {
		return IRCMessage{Raw: raw}
	}

	if raw[0] == ':' {
		// Prefix exists
		parts := strings.SplitN(raw, " :", 2)
		prefix := parts[0]
	}
}

func ExtractPingPongCode(s string) string {
	parts := strings.Split(strings.Trim(s, "\x00\r\n"), " ")
	code := strings.Split(parts[1], "\n")[0]

	return code
}
