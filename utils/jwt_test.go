package utils

import (
	"testing"
)

var secret = "howareyoutoday"

func TestToken(t *testing.T) {
	type args struct {
		Obj map[string]interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{"TestToken", args{
			Obj: map[string]interface{}{
				"username":    "alice",
				"password":    "123456",
				"expire_time": CurrentTime(),
			},
		},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := CreateToken(tt.args.Obj, secret)
			if err != nil {
				t.Errorf("CreateToken returns error: %v.", err)
			} else {
				//
			}

			userName, password, expireTime, _, err := ParseToken(token, secret)
			if err != nil {
				t.Errorf("ParseToken returns error: %v.", err)
			} else {
				if userName != tt.args.Obj["username"] {
					t.Errorf("%s fail, want: %s, got: %s.", tt.name, tt.args.Obj["username"], userName)
				} else if password != tt.args.Obj["password"] {
					t.Errorf("%s fail, want: %s, got: %s.", tt.name, tt.args.Obj["password"], password)
				} else if expireTime <= 0 {
					t.Errorf("%s fail, expire time is invalid: %v.", tt.name, expireTime)
				} else {
					//
				}
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	tests := []struct {
		Name  string
		token string
	}{
		{"TestParseToken", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVfdGltZSI6MTY1OTQ5OTM3MSwicGFzc3dvcmQiOiIxMjM0NTYiLCJ1c2VybmFtZSI6ImFsaWNlIn0.ik1FWNUDFVIZ3yrgD-D0VXkF3mtNkrgLRNH17Mxap04"},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			userName, password, expireTime, _, err := ParseToken(tt.token, secret)
			if err != nil {
				t.Errorf("%s failed, error: %v.", tt.Name, err.Error())
			} else {
				if len(userName) > 0 && len(password) > 0 {
					t.Logf("username: %s, password: %s, expire time: %d.", userName, password, expireTime)
				} else {
					t.Errorf("%s failed: no username or password found.", tt.Name)
				}
			}
		})
	}
}
