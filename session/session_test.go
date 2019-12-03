package session

import (
	"encoding/json"
	"testing"
	"time"
)

var (
	// 默认值:RSA token钥匙对
	testTokenRSAPrivateKeyPEM = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpAIBAAKCAQEAsjmV8ibV3UrQQATxxomfoN20H00B0JyOVYPxRa86g6+G80wv
8NzJApYJp3nQ6266TzhiPtmqHHBsZ7NPZIVyU/rmM04xTb+5PKzp8c7ZfkN8fCs5
OeNSq0+QuxMNW/zwb3TwMzEcHQF42SvP/9PMdN5bOoDWzwRf6ohzE4mPGhLgUQen
jGmX3U4dxBB5smMVynSWNBkF7Ox9hFWT0vlcNDasC9krU3v7S2jOiTT6cF/53+cS
+LfJdNHURlQIna+PfTET8vAndfL9J7RZG/TLk9oYi5oH6Q/1YVVLWNbvlfNEwp9F
POpWhBJ8VpsskbD2+dz6JQ0040YZq80S/sR91QIDAQABAoIBAAzaJjv+HPIGURos
wRqmFgLNug1/yh+3CnSUPTPfPQL+B2dIGTTvVvd+xldza2Nu6rSHxl3t6FyApvCb
d6AyF1qC/1K30spehwcfGQe8+OYoC7QQkeHgyLdd13wGFbKKfPQspJ2sbvCQiJxw
kmFKbrYGyuFfJR5snFYXXOUNyGoCz6UMAENv0YGy9Tz4+kk335QdTG8pcQLQI8vd
ddbbxiKTjQBBqndBmZ8eeIZvd8cjNlT5nlNd3t/GLntPG5g22XLqoH23xB9C1Ikx
qxFDZaQeiEkYMemH1n8PTyaxVP0pq6Rn9x2FvK+6A7ODDjm9A+6Hv5FGxlGpReQS
0KO6oEkCgYEA2zsm8CGpzTmB0S2I6K/0UxiOONFl+PlR9D8E/O0V/Z+9IE6QU9M0
kb31Qpb2VoKyDuHKmzm0AN1RLOG91U36JZVlQcsttqzDFr4N0O7I9Q93KAd/S0z4
9dsDFRP4YyzmMGBgfYHCx+8seK7axVMpJgnnN0bW/MWzez5cgOF6gosCgYEA0B3O
Om08i2WyI6PznlMG0oNEetQGyR5xVc98DG0dDyZr6xYVquy4d+wh/QA4dB58VBG5
DWoRa8QwQQrMNiozNgZX0wDaSImKsoSdhpv5hUpqt8SJ8bWrvkHc3yqEbfDeqvfF
eAXcXCAuhcsjAQoi+PAxqouIQIrMYvxn+Ti87R8CgYEAywKrvBJwOyrVm+6eqVrG
1WwXx2WhGD1INvVkmRKzGnmhmRknbXhXZd6SD2fcFaBRYpaUF8oHdgV79iPUtoHO
8p61dYfAfTjeL2EvShrU3JnFrbvDlOdiY8i7wfkMOkqJnqKgt5hB1wMUG21QCQpJ
QIBLLFTdIJWy7p2A65fg6qECgYEAi0bXJAzEzvlQ/T8Uo6km0K0eoCDTJbdk26uO
dfZz0xbAdESEXa4sSb1ShbGnFjbst4pg0JRicj+Kl4y1W65kNUyLa9+PNaoukwfj
MBmkJErIHpG+S29sL1h+iy82Dyl6qupEUe2CKnkzCUEH/QMeooCEjIDyv1bkL36B
yqzo6rsCgYBs0D18AKdkhlBgdhXJ1KIEzzFNvxZ5u9aa85y9L32dlgWJj84ZMQ2O
GOUD2/ySUXL30QXqwly7qUBLV3dtcKXS0imOelySfQJ+QoZTbrzDI66jbQhJFx+y
8PPBsLv/X+GynNymIxeaelSYcjuHjg40amgA/IlkIOdu4/+tetjUYw==
-----END RSA PRIVATE KEY-----`)
	testTokenRSAPublicKeyPEM = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsjmV8ibV3UrQQATxxomf
oN20H00B0JyOVYPxRa86g6+G80wv8NzJApYJp3nQ6266TzhiPtmqHHBsZ7NPZIVy
U/rmM04xTb+5PKzp8c7ZfkN8fCs5OeNSq0+QuxMNW/zwb3TwMzEcHQF42SvP/9PM
dN5bOoDWzwRf6ohzE4mPGhLgUQenjGmX3U4dxBB5smMVynSWNBkF7Ox9hFWT0vlc
NDasC9krU3v7S2jOiTT6cF/53+cS+LfJdNHURlQIna+PfTET8vAndfL9J7RZG/TL
k9oYi5oH6Q/1YVVLWNbvlfNEwp9FPOpWhBJ8VpsskbD2+dz6JQ0040YZq80S/sR9
1QIDAQAB
-----END PUBLIC KEY-----`)
	// 默认值:HMAC密钥
	testTokenHmacSecret = []byte("sample hmac secret key for token")
)

// 测试通用生成
func Test_SessionParse(t *testing.T) {
	var (
		testMap = map[string]interface{}{
			"foo": "bar",
			"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
			"exp": time.Now().Add(time.Second * 2).Unix(),
		}
		tokenString string
		err         error
	)

	// rsa
	if tokenString, err = SessionEncode(testMap, EncodeRsa); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("rsa token: %s", tokenString)
	}

	//time.Sleep(time.Second * 3) // debug

	if _m, _err := SessionDecode(tokenString); _err != nil {
		t.Fatal(_err)
	} else {
		b, _ := json.Marshal(_m)
		t.Logf("get json: %s", string(b))
	}

	// hmac
	if tokenString, err = SessionEncode(testMap, EncodeHmac); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("hmac token: %s", tokenString)
	}
	if _m, _err := SessionDecode(tokenString); _err != nil {
		t.Fatal(_err)
	} else {
		b, _ := json.Marshal(_m)
		t.Logf("get json: %s", string(b))
	}
}

// 测试session特定格式
func Test_SessionParseAuto(t *testing.T) {
	if s, err := SessionEncodeAuto(&Session{
		UID:   "11111111-1111-1111-1111-111111111111",
		Level: SessionLevelSecure,
	}); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("get token: %s", s)
		// parse
		if s2, err := SessionDecodeAuto(s); err != nil {
			t.Fatal(err)
		} else {
			b, _ := json.Marshal(s2)
			t.Logf("get json: %s", string(b))
		}
	}
}

// 测试token实例
func Test_TokenInstance(t *testing.T) {
	var (
		token    = new(Key)
		tokenStr string
		err      error
		val      = map[string]interface{}{
			"scope":  "nim",
			"acc":    "11111111111111111111111111111111",
			"create": time.Now().Unix(),
		}
	)
	token.DefaultMethod = EncodeRsa
	token.RSAPrivateKeyPEM = testTokenRSAPrivateKeyPEM
	token.RSAPublicKeyPEM = testTokenRSAPublicKeyPEM
	if err = token.Init(); err != nil {
		t.Fatal(err)
	}
	// 生成一个token
	if tokenStr, err = token.TokenEncodeAuto(val); err != nil {
		t.Fatal(err)
	} else {
		t.Logf("token len(%d) %s", len(tokenStr), tokenStr)
	}
	// 解析
	if _m, _err := token.TokenDecode(tokenStr); _err != nil {
		t.Fatal(_err)
	} else {
		b, _ := json.Marshal(_m)
		t.Logf("get %s", string(b))
	}
	// 取sha256
	if _s, _err := token.TokenHashSha256(val); _err != nil {
		t.Fatal(_err)
	} else {
		t.Logf("hash len(%d) %s", len(_s), _s)
	}
}
