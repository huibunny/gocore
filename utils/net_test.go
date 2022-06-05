package utils

import "testing"

func TestGetHostIP(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
		// change to your own local ip address.
		{"TestGetHostIP", "192.168.2.148"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetHostIP(); got != tt.want {
				t.Errorf("GetHostIP() = %v, want %v", got, tt.want)
			}
		})
	}
}
