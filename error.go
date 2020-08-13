package contrib

import "fmt"

var (
	ErrSessionKeyUndefined            = &Error{1, "session key undefined", "contrib"}         // session key undefined
	ErrSessionUndefined               = &Error{2, "session undefined", "contrib"}             // session undefined
	ErrSessionExpUndefined            = &Error{3, "session expire time undefined", "contrib"} // session expire time undefined
	ErrSessionUidUndefined            = &Error{4, "session uid undefined", "contrib"}         // 没有定义uid的session视为非法
	ErrSessionLevelInvalid            = &Error{5, "session level invalid", "contrib"}         // session会话级别非法
	ErrSessionExpMaxOutOfRange        = &Error{6, "session exp out of range", "contrib"}      // session会话有效期超出范围
	ErrTokenParseTypeNotSupport       = &Error{7, "session type not support", "contrib"}      // 不支持的解析类型
	ErrTokenParseSignMethodNotSupport = &Error{8, "unexpected signing method", "contrib"}     // 不支持的解析方法
	ErrTokenInvalid                   = &Error{9, "token invalid", "contrib"}                 // token.valid为false
	ErrTokenClaimsInvalid             = &Error{10, "token claims invalid", "contrib"}         // token.Claims无法解析出
	ErrTokenChangePassword            = &Error{11, "token change password", "contrib"}        // token密码改变
	ErrTokenChangeFreeze              = &Error{12, "token change freeze", "contrib"}          // token状态改变
	ErrTokenChangeLogout              = &Error{13, "token change logout", "contrib"}          // token状态改变
	ErrTokenReqHeaderOrFormKeyInvalid = &Error{14, "request header key or form key invalid by token", "contrib"}
)

// Error translate
type Error struct {
	code   int32
	detail string
	prefix string
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

func (e *Error) Error() string {
	return e.GetPrefix() + "," + e.GetDetail()
}
