package session

import (
	"testing"
	"time"
)

var (
	methods = []string{"RS256", "ES256", "HS256"}
	se      = &Session{
		ExpireTime: int64(TokenExpNormal),
		CreateTime: time.Now().Unix(),
		UID:        "1222222222222222",
		Password:   "123456",
		Level:      SessionLevelNormal,
	}
)

func TestKey_SessionEncode_SessionDecode(t *testing.T) {
	var (
		err error
	)
	for _, v := range methods {
		switch v {
		case "RS256":
			var token string
			key := NewKey("RS256")
			if err = key.SetRSA("./testfile/rsa256-private.pem", "./testfile/rsa256-public.pem"); err != nil {
				t.Fatal(err)
			}
			if token, err = key.SessionEncode(se); err != nil {
				t.Fatal(err)
			}
			t.Log(token)
			var (
				s *Session
			)
			if s, err = key.SessionDecode(token); err != nil {
				t.Fatal(err)
			}
			t.Log(s)
		case "ES256":
			key := NewKey("ES256")
			token := ""
			if err = key.SetECDSA("./testfile/ec256-private.pem", "./testfile/ec256-public.pem"); err != nil {
				t.Fatal(err)
			}
			if token, err = key.SessionEncode(se); err != nil {
				t.Fatal(err)
			}
			t.Log(token)
			var (
				s *Session
			)
			if s, err = key.SessionDecode(token); err != nil {
				t.Fatal(err)
			}
			t.Log(s)
		case "HS256":
			key := NewKey("HS256")
			token := ""
			if err = key.SetHmac("./testfile/hmacTestKey"); err != nil {
				t.Fatal(err)
			}
			if token, err = key.SessionEncode(se); err != nil {
				t.Fatal(err)
			}
			t.Log(token)
			var (
				s *Session
			)
			if s, err = key.SessionDecode(token); err != nil {
				t.Fatal(err)
			}
			t.Log(s)
		}
	}
}
