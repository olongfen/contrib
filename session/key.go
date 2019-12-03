package session

import (
	middle "github.com/srlemon/go-middle"

	"github.com/dgrijalva/jwt-go"

	"crypto/rsa"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
)

// 密钥实例
type Key struct {
	RSAPrivateKeyPEM []byte // 私钥
	RSAPublicKeyPEM  []byte // 公钥
	HmacSecret       []byte // HMAC密钥
	DefaultMethod    string // 默认加密方法
	// hook
	HookSessionCheck func(session *Session) error  // 二次检测session合法性
	HookTokenCheck   func(token interface{}) error // 二次检测token合法性
	// other
	rsaPrivateKey *rsa.PrivateKey
	rsaPublicKey  *rsa.PublicKey
}

// 解析出session
func (d *Key) SessionDecodeAuto(inf interface{}) (ret *Session, err error) {
	if d == nil {
		err = middle.ErrSessionKeyUndefined
		return
	}
	var (
		val map[string]interface{}
	)
	//if val, err = SessionDecode(inf); err != nil {
	//	return
	//}
	if val, err = d.TokenDecode(inf); err != nil {
		return
	}

	// 手动解析
	ret = new(Session)
	ret.CacheToken = inf
	// 逻辑属性
	if v, ok := val[TokenTagExp]; ok == true {
		switch v_ := v.(type) {
		case int32:
			ret.ExpireTime = int64(v_)
			break
		case int64:
			ret.ExpireTime = v_
			break
		case float64:
			ret.ExpireTime = int64(v_)
			break
		default:
			err = fmt.Errorf("unknown %s value: %v", TokenTagExp, v)
			break
		}
	}
	// 内容
	if v, ok := val[TokenTagCre]; ok == true {
		if s, ok := v.(float64); ok == true {
			ret.CreateTime = int64(s)
		}
	}
	if v, ok := val[TokenTagUid]; ok == true {
		if s, ok := v.(string); ok == true {
			ret.UID = s
		}
	}
	if v, ok := val[TokenTagLevel]; ok == true {
		if s, ok := v.(string); ok == true {
			ret.Level = s
		}
	}
	if v, ok := val[TokenTagPsw]; ok == true {
		if s, ok := v.(string); ok == true {
			ret.Password = s
		}
	}
	if v, ok := val[TokenTagIp]; ok == true {
		if s, ok := v.(string); ok == true {
			ret.IP = s
		}
	}
	if v, ok := val[TokenTagId]; ok == true {
		if s, ok := v.(string); ok == true {
			ret.ID = s
		}
	}

	// 合法性
	if err = ret.Valid(); err != nil {
		return
	}
	if d.HookSessionCheck != nil {
		if err = d.HookSessionCheck(ret); err != nil {
			return
		}
	}
	return
}

// 将session编码为token
func (d *Key) SessionEncodeAuto(s *Session) (token string, err error) {
	if d == nil {
		err = middle.ErrSessionKeyUndefined
		return
	}
	if s == nil {
		err = middle.ErrSessionUndefined
		return
	} else if err = s.Valid(); err != nil {
		return
	}

	// 手动填充map
	var m = make(map[string]interface{})

	// 必填:逻辑属性
	m[TokenTagExp] = s.ExpireTime // ExpireTime

	// 可选内容
	if s.CreateTime > 0 {
		// cre
		m[TokenTagCre] = s.CreateTime
	}
	if len(s.UID) > 0 {
		// uid
		m[TokenTagUid] = s.UID
	}
	if len(s.Level) > 0 {
		// level
		m[TokenTagLevel] = s.Level
	}
	if len(s.Password) > 0 {
		// psw
		m[TokenTagPsw] = s.Password
	}
	if len(s.IP) > 0 {
		// ip
		m[TokenTagIp] = s.IP
	}
	if len(s.ID) > 0 {
		// id
		m[TokenTagId] = s.ID
	}

	token, err = d.TokenEncode(m, d.DefaultMethod) // 默认加密
	return
}

// 解析出需要的值
func (d *Key) TokenDecode(inf interface{}) (ret map[string]interface{}, err error) {
	if d == nil {
		err = middle.ErrSessionKeyUndefined
		return
	}
	switch v := inf.(type) {
	case string:
		ret, err = d.tokenParse(v)
		break
	case []byte:
		ret, err = d.tokenParse(string(v))
		break
	case json.RawMessage:
		ret, err = d.tokenParse(string(v))
		break
	case *http.Request:
		var token string
		if token, err = PubGetTokenFromReq(v); err != nil {
			return
		}
		ret, err = d.tokenParse(token)
		break
	default:
		err = middle.ErrTokenParseTypeNotSupport
		break
	}
	if err == nil && d.HookTokenCheck != nil {
		if err = d.HookTokenCheck(inf); err != nil {
			return
		}
	}
	return
}

// 将值编码为token
func (d *Key) TokenEncode(val map[string]interface{}, method string) (token string, err error) {
	if d == nil {
		err = middle.ErrSessionKeyUndefined
		return
	}
	if val == nil {
		err = middle.ErrTokenClaimsInvalid
		return
	}

	switch method {
	case EncodeRsa:
		token, err = d.tokenEncodeRsa(val)
		break
	case EncodeHmac:
		token, err = d.tokenEncodeHmac(val)
		break
	default:
		err = middle.ErrTokenParseSignMethodNotSupport
		break
	}
	return
}

// 将值编码为token
func (d *Key) TokenEncodeAuto(val map[string]interface{}) (token string, err error) {
	return d.TokenEncode(val, d.DefaultMethod)
}

// 取值的token的sha256
func (d *Key) TokenHashSha256(val map[string]interface{}) (sign string, err error) {
	var (
		token string
	)
	if token, err = d.TokenEncodeAuto(val); err != nil {
		return
	}
	sign = fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
	return
}

// 设置RSA密钥对
func (d *Key) SetRSA(priPem []byte, pubPem []byte) (err error) {
	if d == nil {
		err = middle.ErrSessionKeyUndefined
		return
	}
	d.RSAPrivateKeyPEM = priPem
	d.RSAPublicKeyPEM = pubPem
	if err = d.Init(); err != nil {
		return
	}
	d.DefaultMethod = EncodeRsa
	return
}

// 设置HMAC密钥
func (d *Key) SetHmac(hmacPri []byte) (err error) {
	if d == nil {
		err = middle.ErrSessionKeyUndefined
		return
	}
	d.HmacSecret = hmacPri
	if err = d.Init(); err != nil {
		panic(err)
	}
	d.DefaultMethod = EncodeHmac
	return
}

// 解析出token中的值
func (d *Key) tokenParse(tokenStr string) (ret map[string]interface{}, err error) {
	var (
		token *jwt.Token
	)

	// 解析token
	if token, err = jwt.Parse(tokenStr, d.parseKey); err != nil {
		// 重定向错误
		err = middle.ErrTokenInvalid.SetVars("crypto/rsa")
		return
	} else if token.Valid == false {
		err = middle.ErrTokenInvalid
		return
	}

	// 解析map
	if claims, ok := token.Claims.(jwt.MapClaims); ok == true {
		//// 解析合法性
		//if err = claims.Valid(); err != nil {
		//	return
		//}
		ret = map[string]interface{}(claims)
	} else {
		err = middle.ErrTokenClaimsInvalid
		return
	}

	return
}

// 将map转为rsa的token
func (d *Key) tokenEncodeRsa(val map[string]interface{}) (tokenStr string, err error) {
	tokenStr, err = jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(val)).SignedString(d.rsaPrivateKey)
	return
}

// 将map转为hmac的token
func (d *Key) tokenEncodeHmac(val map[string]interface{}) (tokenStr string, err error) {
	tokenStr, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(val)).SignedString(d.HmacSecret)
	return
}

// 按类型取出密钥
func (d *Key) parseKey(token *jwt.Token) (key interface{}, err error) {
	switch token.Method.(type) {
	case *jwt.SigningMethodRSA:
		// RSA
		key = d.rsaPublicKey
		break
	case *jwt.SigningMethodHMAC:
		// HMAC
		key = d.HmacSecret
		break
	default:
		err = fmt.Errorf("%s '%v'", middle.ErrTokenParseSignMethodNotSupport.Error(), token.Header["alg"])
		break
	}
	return
}

// init
func (d *Key) Init() (err error) {
	if d == nil {
		err = middle.ErrSessionKeyUndefined
		return
	}

	// default method
	if len(d.DefaultMethod) == 0 {
		d.DefaultMethod = CfgDefaultMethod
	}

	// rsa instance/replace
	if len(d.RSAPrivateKeyPEM) > 0 {
		if d.rsaPrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(d.RSAPrivateKeyPEM); err != nil {
			return
		}
	}
	if len(d.RSAPublicKeyPEM) > 0 {
		if d.rsaPublicKey, err = jwt.ParseRSAPublicKeyFromPEM(d.RSAPublicKeyPEM); err != nil {
			return
		}
	}

	//// hmmc: default
	//if len(d.HmacSecret) == 0 {
	//	d.HmacSecret = defaultHmacSecret
	//}

	return
}

// new key
func NewKey(in *Key) (d *Key, err error) {
	if in != nil {
		d = in
	} else {
		d = new(Key)
	}
	err = d.Init()
	return
}
