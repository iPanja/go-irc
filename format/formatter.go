package format

import (
	"strings"
)

// FormatUserInput will perform special transformations on the user's input.
//
// Currently, this only transforms `/<cmd> ...` => `CMD ...`
func FormatUserInput(input string) string {
	// Format command type inputs
	// Ex: "/join #foobar" => "JOIN #foobar"
	if len(input) > 0 && input[0] == '/' {
		parts := strings.Split(input[1:], " ")

		parts[0] = strings.ToUpper(parts[0])
		result := strings.Join(parts, " ")

		input = result
	}

	return input
}
