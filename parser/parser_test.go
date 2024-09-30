package parser

import (
	"reflect"
	"testing"
)

func TestNewIrcMessage(t *testing.T) {
	tests := []struct {
		input string
		want  IRCMessage
	}{
		{
			input: ":ZeroOne!~ZeroOne@freenode-1u5.5hk.gtmbaa.IP PRIVMSG #trivia :10,1[3,1No. 1960 art-and-literature10,1]15,1    Shakespeare: The Abbot Of Westminster Is A Character In ____ _______",
			want: IRCMessage{
				Prefix:  ":ZeroOne!~ZeroOne@freenode-1u5.5hk.gtmbaa.IP",
				Command: "PRIVMSG",
				Params:  []string{"#trivia", ":10,1[3,1No.", "1960", "art-and-literature10,1]15,1", "", "", "", "Shakespeare:", "The", "Abbot", "Of", "Westminster", "Is", "A", "Character", "In", "____", "_______"},
				Raw:     ":ZeroOne!~ZeroOne@freenode-1u5.5hk.gtmbaa.IP PRIVMSG #trivia :10,1[3,1No. 1960 art-and-literature10,1]15,1    Shakespeare: The Abbot Of Westminster Is A Character In ____ _______",
			},
		},
	}

	for _, test := range tests {
		result := NewIRCMessage(test.input)

		if !reflect.DeepEqual(result, test.want) {
			t.Errorf("NewIRCMessage() = %v, want %v", result.Params, test.want.Params)
		}
	}
}
