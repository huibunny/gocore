package utils

import "testing"

func Test_Crypt(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{"Test_Crypt", "123456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashPass, err := HashPassword(tt.password)
			if err != nil {
				t.Errorf("%s failed to hash, error: %v.", tt.name, err)
			} else {
				if CheckPasswordHash(tt.password, hashPass) {
					//
				} else {
					t.Fatalf("%s failed to check, error: %v.", tt.name, err)
				}
			}
		})
	}
}
