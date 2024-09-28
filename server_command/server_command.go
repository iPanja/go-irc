package server_command

import (
	"regexp"
)

type ServerCommand interface {
	// IsSuccessful will check the status of a user-issued command/action based on the server's response.
	//
	// Returns: (was command successful?, data extracted from response)
	IsSuccessful(response string) (bool, map[string]string)
}

// SCJoin | First response from server after initial handshake.
//
// Expected response: :<HOSTNAME> 002 <NICK> :Your host is <HOSTNAME>, running version <VERSION>
//
// Returns: ok, {"hostname": <HOSTNAME>}
type SCJoin struct{}

func (SCJoin) IsSuccessful(response string) (bool, map[string]string) {
	re := regexp.MustCompile("^:(?P<hostname>.+) 002 .+:Your host is")

	matches := re.FindStringSubmatch(response)
	hostIndex := re.SubexpIndex("hostname")

	if len(matches) == 0 || hostIndex == -1 {
		return false, map[string]string{}
	}

	return true, map[string]string{"hostname": matches[hostIndex]}
}
