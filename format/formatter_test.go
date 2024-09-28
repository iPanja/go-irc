package format

import (
	"testing"
)

func TestFormatUserInput(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{
			input: "/join #foobar",
			want:  "JOIN #foobar",
		},
		{
			input: "test message",
			want:  "test message",
		},
		{
			input: "!server_cmd",
			want:  "!server_cmd",
		},
	}

	for _, test := range tests {
		result := FormatUserInput(test.input)

		if result != test.want {
			t.Errorf("formatUserInput(%q) = %q, want %q", test.input, result, test.want)
		}
	}
}
