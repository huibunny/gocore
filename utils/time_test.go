package utils

import (
	"testing"
)

func Test_CurrentTime(t *testing.T) {
	tests := []struct {
		Name string
	}{
		{"TestCurrentTime"},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			now := CurrentTime()
			t.Logf("now: %v.", now)
		})
	}
}
