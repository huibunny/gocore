package consul

import (
	"testing"
)

func Test_RetrieveAddressPort(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		// TODO: Add test cases.
		{"test1", args{"HTTP GET http://192.168.2.180:9898/app/healthcheck: 200 OK Output: {\"status\": true}"}, "192.168.2.180", "9898"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := RetrieveAddressPort(tt.args.url)
			if got != tt.want {
				t.Errorf("RetrieveAddressPort() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("RetrieveAddressPort() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestRetrieveServiceName(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"testRetrieveServiceName", args{"Service 'helloapp' check"}, "helloapp"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RetrieveServiceName(tt.args.url); got != tt.want {
				t.Errorf("RetrieveServiceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRetrieveServiceID(t *testing.T) {
	type args struct {
		checkID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{"testRetrieveServiceID", args{"service:helloapp-3c17f434b5c34311a068d710624c308a"}, "helloapp-3c17f434b5c34311a068d710624c308a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RetrieveServiceID(tt.args.checkID); got != tt.want {
				t.Errorf("RetrieveServiceID() = %v, want %v", got, tt.want)
			}
		})
	}
}
