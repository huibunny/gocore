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
				"user_id":     "10337fc7-a6f1-4343-9561-df56df9011b0",
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

			userId, expireTime, _, err := ParseToken(token, secret)
			if err != nil {
				t.Errorf("ParseToken returns error: %v.", err)
			} else {
				if userId != tt.args.Obj["user_id"] {
					t.Errorf("%s fail, want: %s, got: %s.", tt.name, tt.args.Obj["user_id"], userId)
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
		{"TestParseToken", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmVfdGltZSI6MTY1OTUzMDYxMCwidXNlcl9pZCI6IjEwMzM3ZmM3LWE2ZjEtNDM0My05NTYxLWRmNTZkZjkwMTFiMCJ9.515r399vCi4GrRGywnU2BBlSLgoMozxqfKCP-Safjs0"},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			userId, expireTime, _, err := ParseToken(tt.token, secret)
			if err != nil {
				t.Errorf("%s failed, error: %v.", tt.Name, err.Error())
			} else {
				if len(userId) > 0 {
					t.Logf("user id: %s, expire time: %d.", userId, expireTime)
				} else {
					t.Errorf("%s failed: no user id found.", tt.Name)
				}
			}
		})
	}
}
