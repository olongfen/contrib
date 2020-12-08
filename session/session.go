package session

import (
	error2 "github.com/olongfen/contrib/utils"

	"time"
)

const (
	// 会话级别
	LevelNormal = ""       // 默认会话
	LevelSecure = "secure" // 安全会话
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
	TokenExpNormal = time.Hour * 24 // 默认登录有效期一天
	// 会话最大有效期
	ExpMaxNormal = time.Hour * 24 * 7 // 普通会话最长一周
	ExpMaxSecure = time.Hour * 24     // 安全会话最长一天
)

const (
	// jwt https://tools.ietf.org/html/draft-ietf-oauth-json-web-token-32
	TokenTagIss = "iss"
	TokenTagSub = "sub"
	TokenTagAud = "aud"
	TokenTagExp = "exp" // 超时
	// 自定义部分
	TokenTagUid     = "uid" //
	TokenTagContent = "content"
)

type Params struct {
	EncryptMethod string
	// 逻辑属性
	ExpireTime int64
	UID        string
	Content    *Key
}

// session规范
type Session struct {
	// 逻辑属性
	ExpireTime int64  `json:"expireTime,omitempty"` // 超时时间戳
	UID        string `json:"uid,omitempty"`        // 用户uid,唯一uid
	// 内容
	Content map[string]interface{} `json:"content,omitempty"`
}

// **
func (s *Session) Valid() (err error) {

	// 没有uid的session无价值
	if len(s.UID) == 0 {
		err = error2.ErrSessionUidUndefined
		return
	}

	// expire: normal
	if s.ExpireTime <= 0 {
		err = error2.ErrSessionExpUndefined
		return
	}

	// expire: 以now时间点验证
	var now_ = time.Now()
	level, _ok := s.Content["level"]
	if !_ok {
		level = LevelNormal
	}
	switch level {
	case LevelNormal:
		if s.ExpireTime > now_.Add(ExpMaxNormal).UnixNano() {
			err = error2.ErrSessionExpMaxOutOfRange
			return
		}
		break
	case LevelSecure:
		if s.ExpireTime > now_.Add(ExpMaxSecure).UnixNano() {
			err = error2.ErrSessionExpMaxOutOfRange
			return
		}
		break
	}

	return
}
