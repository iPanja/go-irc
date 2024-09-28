package parser

import (
	"reflect"
	"testing"
)

func TestSCJoin_IsSuccessful(t *testing.T) {
	tests := []struct {
		input    string
		wantOk   bool
		wantData map[string]string
	}{
		{
			input:    ":*.freenode.net 002 <NICK> :Your host is *.freenode.net, running version InspIRCd-3",
			wantOk:   true,
			wantData: map[string]string{"hostname": "*.freenode.net"},
		},
		{
			input:    ":some random response!! HEHEHE",
			wantOk:   false,
			wantData: map[string]string{},
		},
	}

	sc := SCJoin{}
	for _, test := range tests {
		ok, data := sc.IsSuccessful(test.input)

		if ok != test.wantOk {
			t.Errorf("TestSCJoin_IsSuccessful: (ok) got %v want %v", ok, test.wantOk)
		}

		if !reflect.DeepEqual(data, test.wantData) {
			t.Errorf("TestSCJoin_IsSuccessful: (data) got %v want %v", data, test.wantData)
		}
	}
}
