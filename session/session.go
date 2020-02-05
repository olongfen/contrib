package session

import (
	project "github.com/olongfen/contrib"

	"time"
)

const (
	EncodeRsa  = "rsa"
	EncodeHmac = "hmac"
	// 会话级别
	SessionLevelNormal = ""       // 默认会话
	SessionLevelSecure = "secure" // 安全会话
)

const (
	// token验证标记
	TokenCodePass   int8 = iota // 0: 合法
	TokenCodePsw                // 1: 密码已修改
	TokenCodeFreeze             // 2: 账号已冻结
	TokenCodeLogout             // 3: 用户已在token过期前安全登出
)

var (
	// 有效期
	TokenExpNormal = time.Hour * 24     // 默认登录有效期一天
	TokenExpLong   = time.Hour * 24 * 7 // 超长登录一周
	TokenExpSecure = time.Minute * 30   // 默认安全会话30分钟
	// 会话最大有效期
	SessionExpMaxNormal = time.Hour * 24 * 365 // 普通会话最长一年
	SessionExpMaxLong   = time.Hour * 24 * 365 // 超长会话最长一年
	SessionExpMaxSecure = time.Hour * 24       // 安全会话最长一天
	//
	CfgDefaultMethod = EncodeRsa // 默认加密方法
	KeyDefault       *Key        // 默认key实例
)

const (
	// jwt https://tools.ietf.org/html/draft-ietf-oauth-json-web-token-32
	TokenTagIss = "iss"
	TokenTagSub = "sub"
	TokenTagAud = "aud"
	TokenTagExp = "exp" // 超时
	TokenTagNbf = "nbf" // 在此之前无效
	TokenTagIat = "iat"
	TokenTagJti = "jti"
	// 自定义部分
	TokenTagCre   = "cre"   // 创建于
	TokenTagUid   = "uid"   //
	TokenTagLevel = "level" //
	TokenTagPsw   = "psw"   //
	TokenTagIp    = "ip"    //
	TokenTagId    = "id"    //
)

// session规范
type Session struct {
	// 逻辑属性
	ExpireTime int64 `json:"expireTime,omitempty"` // 超时时间戳
	// 内容
	CreateTime int64  `json:"createTime,omitempty"` // 创建时间戳
	UID        string `json:"uid,omitempty"`        // 用户uid
	Level      string `json:"level,omitempty"`      // 会话等级
	IP         string `json:"ip,omitempty"`         // 登陆地址
	ID         string `json:"id,omitempty"`         // 唯一标记/设备id
	// 验证
	Password string `json:"password,omitempty"` // 密码哈希摘要
	// cache
	CacheToken interface{} `json:"-"` //
}

// 解析出session
func SessionDecodeAuto(inf interface{}) (ret *Session, err error) {
	if KeyDefault == nil {
		err = project.ErrSessionKeyDefaultUndefined
		return
	}
	ret, err = KeyDefault.SessionDecodeAuto(inf)
	return
}

// 将session编码为token
func SessionEncodeAuto(s *Session) (token string, err error) {
	if KeyDefault == nil {
		err = project.ErrSessionKeyDefaultUndefined
		return
	}
	token, err = KeyDefault.SessionEncodeAuto(s)
	return
}

// 解析出需要的值
func SessionDecode(inf interface{}) (ret map[string]interface{}, err error) {
	if KeyDefault == nil {
		err = project.ErrSessionKeyDefaultUndefined
		return
	}
	ret, err = KeyDefault.TokenDecode(inf)
	return
}

// 将值编码为token
func SessionEncode(val map[string]interface{}, method string) (token string, err error) {
	if KeyDefault == nil {
		err = project.ErrSessionKeyDefaultUndefined
		return
	}
	token, err = KeyDefault.TokenEncode(val, method)
	return
}

// 设置RSA密钥对
func SessionSetRSA(priPem []byte, pubPem []byte) (err error) {
	if KeyDefault == nil {
		err = project.ErrSessionKeyDefaultUndefined
		return
	}
	if err = KeyDefault.SetRSA(priPem, pubPem); err != nil {
		panic(err)
	}
	return
}

// 设置HMAC密钥
func SessionSetHmac(hmacPri []byte) (err error) {
	if KeyDefault == nil {
		err = project.ErrSessionKeyDefaultUndefined
		return
	}
	if err = KeyDefault.SetHmac(hmacPri); err != nil {
		panic(err)
	}
	return
}

// **
func (s *Session) Valid() (err error) {

	// 没有uid的session无价值
	if len(s.UID) == 0 {
		err = project.ErrSessionUidUndefined
		return
	}

	// expire: normal
	if s.ExpireTime <= 0 {
		err = project.ErrSessionExpUndefined
		return
	}

	// expire: 以now时间点验证
	var now_ = time.Now()
	switch s.Level {
	case SessionLevelNormal:
		if s.ExpireTime > now_.Add(SessionExpMaxNormal).Unix() {
			err = project.ErrSessionExpMaxOutOfRange
			return
		}
		break
	case SessionLevelSecure:
		if s.ExpireTime > now_.Add(SessionExpMaxSecure).Unix() {
			err = project.ErrSessionExpMaxOutOfRange
			return
		}
		break
	default:
		err = project.ErrSessionLevelInvalid
		return
	}

	return
}

func init() {
	// key
	if err := InitKey(); err != nil {
		panic(err)
	}
}
