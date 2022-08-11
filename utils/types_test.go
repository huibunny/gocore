package utils

import (
	"testing"
)

func Test_ToStrings(t *testing.T) {
	tests := []struct {
		Name   string
		values []interface{}
	}{
		{"TestCurrentTime", []interface{}{"hello", "world"}},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			strValues := ToStrings(tt.values)
			t.Logf("%v.", strValues)
		})
	}
}
