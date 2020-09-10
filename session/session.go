package session

import (
	project "github.com/olongfen/contrib"

	"time"
)

const (
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
	TokenExpTenMinute = time.Minute * 10
	TokenExpHour      = time.Hour
	TokenExpNormal    = time.Hour * 24     // 默认登录有效期一天
	TokenExpLong      = time.Hour * 24 * 7 // 超长登录一周
	TokenExpSecure    = time.Minute * 30   // 默认安全会话30分钟
	// 会话最大有效期
	SessionExpMaxNormal = time.Hour * 24 * 365 // 普通会话最长一年
	SessionExpMaxLong   = time.Hour * 24 * 365 // 超长会话最长一年
	SessionExpMaxSecure = time.Hour * 24       // 安全会话最长一天
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
	TokenTagCre      = "cre"      // 创建于
	TokenTagUid      = "uid"      //
	TokenTagLevel    = "level"    //
	TokenTagPsw      = "psw"      //
	TokenTagIp       = "ip"       //
	TokenTagId       = "id"       //
	TokenTagUsername = "username" //
	TokenTagDeviceId = "deviceId"
)

type Params struct {
	EncryptMethod string
	// 逻辑属性
	ExpireTime int64
	UID        string
	Level      string
	IP         string
	ID         int64
	DeviceID   string
	Username   string
	// 验证
	Password string
}

// session规范
type Session struct {
	// 逻辑属性
	ExpireTime int64 `json:"expireTime,omitempty"` // 超时时间戳
	// 内容
	CreateTime int64  `json:"createTime,omitempty"` // 创建时间戳
	UID        string `json:"uid,omitempty"`        // 用户uid,唯一uid
	Level      string `json:"level,omitempty"`      // 会话等级
	IP         string `json:"ip,omitempty"`         // 登陆地址
	ID         int64  `json:"id,omitempty"`         // 唯一标记,用户id主键
	DeviceID   string `json:"deviceId"`
	Username   string `json:"username"` // 用户名，唯一

	// 验证
	Password string `json:"password,omitempty"` // 密码哈希摘要
	// cache
	CacheToken interface{} `json:"-"` //
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
		if s.ExpireTime > now_.Add(SessionExpMaxNormal).UnixNano() {
			err = project.ErrSessionExpMaxOutOfRange
			return
		}
		break
	case SessionLevelSecure:
		if s.ExpireTime > now_.Add(SessionExpMaxSecure).UnixNano() {
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
