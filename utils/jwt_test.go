package utils

import (
	"testing"
)

func TestToken(t *testing.T) {
	type args struct {
		Obj    map[string]interface{}
		Secret string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"TestToken", args{
			Obj: map[string]interface{}{
				"username": "alice",
				"password": "123456",
			},
			Secret: "howareyoutoday?",
		},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := CreateToken(tt.args.Obj, tt.args.Secret)
			if err != nil {
				t.Errorf("CreateToken returns error: %v.", err)
			} else {
				//
			}

			userInfo, err := ParseToken(token, tt.args.Secret)
			if err != nil {
				t.Errorf("ParseToken returns error: %v.", err)
			} else {
				if userInfo["username"] != tt.args.Obj["username"] {
					t.Errorf("ParseToken fail, want: %s, got: %s.", tt.args.Obj["username"], userInfo["username"])
				} else if userInfo["password"] != tt.args.Obj["password"] {
					t.Errorf("ParseToken fail, want: %s, got: %s.", tt.args.Obj["password"], userInfo["password"])
				} else {
					//
				}
			}
		})
	}
}
