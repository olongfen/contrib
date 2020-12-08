package utils

import (
	"fmt"
)

var (
	ErrSessionKeyUndefined            = &Error{code: 1, detail: "session key undefined", prefix: "contrib"}         // session key undefined
	ErrSessionUndefined               = &Error{code: 2, detail: "session undefined", prefix: "contrib"}             // session undefined
	ErrSessionExpUndefined            = &Error{code: 3, detail: "session expire time undefined", prefix: "contrib"} // session expire time undefined
	ErrSessionUidUndefined            = &Error{code: 4, detail: "session uid undefined", prefix: "contrib"}         // 没有定义uid的session视为非法
	ErrSessionLevelInvalid            = &Error{code: 5, detail: "session level invalid", prefix: "contrib"}         // session会话级别非法
	ErrSessionExpMaxOutOfRange        = &Error{code: 6, detail: "session exp out of range", prefix: "contrib"}      // session会话有效期超出范围
	ErrTokenParseTypeNotSupport       = &Error{code: 7, detail: "session type not support", prefix: "contrib"}      // 不支持的解析类型
	ErrTokenParseSignMethodNotSupport = &Error{code: 8, detail: "unexpected signing method", prefix: "contrib"}     // 不支持的解析方法
	ErrTokenInvalid                   = &Error{code: 9, detail: "token invalid", prefix: "contrib"}                 // token.valid为false
	ErrTokenClaimsInvalid             = &Error{code: 10, detail: "token claims invalid", prefix: "contrib"}         // token.Claims无法解析出
	ErrTokenChangePassword            = &Error{code: 11, detail: "token change password", prefix: "contrib"}        // token密码改变
	ErrTokenChangeFreeze              = &Error{code: 12, detail: "token change freeze", prefix: "contrib"}          // token状态改变
	ErrTokenChangeLogout              = &Error{code: 13, detail: "token change logout", prefix: "contrib"}          // token状态改变
	ErrTokenReqHeaderOrFormKeyInvalid = &Error{code: 14, detail: "request header key or form key invalid by token", prefix: "contrib"}
)

// Error translate
type Error struct {
	code   int32
	detail string
	prefix string
	meta   interface{}
}

// NewError
func NewError(code int32, detail string, prefix string) (ret *Error) {
	ret = new(Error)
	ret.code = code
	ret.detail = detail
	ret.prefix = prefix
	return
}

func (e *Error) GetPrefix() string {
	return e.prefix
}

func (e *Error) GetCode() int32 {
	return e.code
}

func (e *Error) GetKey() string {
	return fmt.Sprintf("prefix:%s , code:%d", e.GetPrefix(), e.GetCode())
}

func (e *Error) GetDetail() string {
	return fmt.Sprintf("detail: %s", e.detail)
}

func (e *Error) SetDetail(c string) *Error {
	e.detail = c
	return e
}

func (e *Error) SetMeta(m interface{}) *Error {
	e.meta = m
	return e
}

func (e *Error) GetMeta() string {
	switch e.meta.(type) {
	case string:
		return fmt.Sprintf("meta: %s", e.meta)
	case int, int32, int16, int8, int64:
		return fmt.Sprintf("meta: %v", e.meta)
	case uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("meta: %v", e.meta)
	case float32, float64:
		return fmt.Sprintf("meta: %v", e.meta)
	default:
		return fmt.Sprintf("meta: %s", JSONMarshalMust(e.meta))
	}
}

func (e *Error) Error() string {
	if e.meta != nil {
		// reset meta
		defer func() {
			e.meta = nil
		}()
		return e.GetPrefix() + "," + e.GetDetail() + "," + e.GetMeta()
	}
	return e.GetPrefix() + "," + e.GetDetail()

}
