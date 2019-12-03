package middle

import "fmt"

var (
	ErrParamInvalid     = &Error{1, "param invalid", "middle", nil}
	ErrInputDataIsError = &Error{2, "Input get information error", "middle", nil}
	ErrWrongPassword=&Error{3,"wrong password","middle",nil}
	ErrPasswordNotSame = &Error{4,"password is not the same","middle",nil}
	ErrTokenUndefined                 = &Error{5, "token undefined", "middle", nil}                                               // token未定义
	ErrCtxUndefined                   = &Error{6, "params context undefined", "middle", nil}                                      // ctx未定义
	ErrCtxCtrlUndefined               = &Error{7 , "ctrl undefined in context", "middle", nil}                                     // ctrl未取得
	ErrCtxReqUndefined                = &Error{8, "request undefined in context", "middle", nil}                                  // req未取得
	ErrPasswordUndefined              = &Error{9, "password undefined", "middle", nil}                                            // 密码未定义
	ErrPasswordPayRequire             = &Error{10, "password-pay require", "middle", nil}                                          // 需要支付密码
	ErrPasswordPayInvalid             = &Error{11, "payment password error", "middle", nil}                                        // 支付密码错误
	ErrCtrlReqJsonUndefined           = &Error{12, "json: can not find string in request.Body", "middle", nil}                     // 无法解析json
	ErrGraphqlMetaUndefined           = &Error{13, "meta undefined", "middle", nil}                                                // meta对象nil
	ErrGraphqlQueryLimitMax           = &Error{14, "query limit max 100", "middle", nil}                                           // limit超出范围
	ErrHookPasswordVerify             = &Error{15, "HookPasswordVerify undefined", "middle", nil}                                  // 密码二次验证函数未定义
	ErrQueryUndefined                 = &Error{16, "query undefined", "middle", nil}                                               // query未定义
	ErrConfigSavePathUndefined        = &Error{17, "config save path undefined", "middle", nil}                                    // 配置文件保存地址未定义
	ErrConfigSavePointUndefined       = &Error{19, "config object undefined", "middle", nil}                                       // 配置文件保存地址未定义
	ErrRouterHandleType               = &Error{20, "handle invalid, allow http.HandlerFunc/*http.HandlerFunc only", "middle", nil} // handle函数非法
	ErrOptionsNotAllow                = &Error{21, "options method not allow", "middle", nil}                                      // 参数不允许
	ErrObjectUndefined                = &Error{22, "object undefined", "middle", nil}                                              // 参数未定义
	ErrPermUndefined                  = &Error{23, "permission undefined", "middle", nil}                                          // 权限未定义
	ErrPermNotAllow                   = &Error{24, "action not allow {{.var1}}", "middle", nil}                                    // 动作不允许 {{.var1}}
	ErrParamUndefined                 = &Error{25, "params undefined", "middle", nil}                                              // 参数未定义
	ErrUndefined                      = &Error{26, "undefined {{.var1}}", "middle", nil}                                           // 未定义 {{.var1}}
	ErrRpcDevelopmentPhaseOnly        = &Error{27, "this rpc can user in development phase only", "middle", nil}                   // 只能在开发阶段
	ErrRpcConnectionRefused           = &Error{28, "rpc connection refused", "middle", nil}                                        // rpc无法连接
	ErrRpcTimeFormatInvalid           = &Error{29, "rpc time format invalid", "middle", nil}                                       // 时间格式错误
	ErrSessionKeyUndefined            = &Error{30, "session key undefined", "middle", nil}                                         // session key undefined
	ErrSessionKeyDefaultUndefined     = &Error{31, "default key instance undefined", "middle", nil}                                // 默认key实例未定义
	ErrSessionUndefined               = &Error{33, "session undefined", "middle", nil}                                             // session undefined
	ErrSessionExpUndefined            = &Error{34, "session expire time undefined", "middle", nil}                                 // session expire time undefined
	ErrSessionUidUndefined            = &Error{35, "session uid undefine", "middle", nil}                                          // 没有定义uid的session视为非法
	ErrSessionLevelInvalid            = &Error{36, "session level invalid", "middle", nil}                                         // session会话级别非法
	ErrSessionExpMaxOutOfRange        = &Error{37, "session exp out of range", "middle", nil}                                      // session会话有效期超出范围
	ErrTokenParseTypeNotSupport       = &Error{38, "session type not support", "middle", nil}                                      // 不支持的解析类型
	ErrTokenParseSignMethodNotSupport = &Error{39, "unexpected signing method", "middle", nil}                                     // 不支持的解析方法
	ErrTokenInvalid                   = &Error{40, "token invalid", "middle", nil}                                                 // token.valid为false
	ErrTokenClaimsInvalid             = &Error{41, "token claims invalid", "middle", nil}                                          // token.Claims无法解析出
	ErrTokenChangePassword            = &Error{42, "token change password", "middle", nil}                                         // token密码改变
	ErrTokenChangeFreeze              = &Error{43, "token change freeze", "middle", nil}                                           // token状态改变
	ErrTokenChangeLogout              = &Error{44, "token change logout", "middle", nil}                                           // token状态改变
	ErrHookPluginNotFound             = &Error{45, "not plugin file(s) found", "middle", nil}                                      // 读不到任何插件文件
	ErrTimeout                        = &Error{46, "timeout {{.var1}}", "middle", nil}                                             // 超时{{.var1}}
	ErrPanicRecover                   = &Error{47, "panic {{.var1}}", "middle", nil}                                               // 恐慌错误 {{.var1}}
	ErrModelInited                    = &Error{48, "model is not init", "middle", nil}                                             // 初始化错误
	ErrOrmMatchNone                   = &Error{49, "match none", "middle", nil}                                                   // 数据库找不到对象
	ErrOrmMatchMultiple               = &Error{50, "match multiple", "middle", nil}                                               // 期望搜到一条记录，但是返回多条
	// ErrParamInvalid                   = &Error{129, "params invalid {{.var1}}", "middle", nil}                                     // 参数非法 {{.var1}}
	ErrSessionExpired                 = &Error{52, "session expired", "middle", nil}                                              // 会话已过期
	ErrTokenExpired                   = &Error{52, "token expired", "middle", nil}                                                // 身份令牌已过期
	ErrPlaceHolder                    = &Error{53, "error placeholder", "middle", nil}                                            // 错误的错误
	ErrConfigInvalid                  = &Error{54, "config file empty or invalid", "middle", nil}                                 // 配置文件非法或不完整
	ErrRbacInherCircle                = &Error{55, "RBAC roles Inher circle {{.var1}} {{.var2}}", "middle", nil}                  // 角色继承循环 {{.var1}} {{.var2}}
	ErrUuidFormat                     = &Error{56, "uuid format invalid {{.var1}}", "middle", nil}                                // uuid 格式错误 {{.var1}}
	ErrApiKeyInvalid                  = &Error{57, "API KEY invalid", "middle", nil}                                              // 接口密钥错误
	ErrLengthOutOfRange               = &Error{58, "length out of range {{.var1}}", "middle", nil}                                // 长度超出范围 {{.var1}}
	ErrParamIllegal                   = &Error{59, "param illegal", "middle", nil}
	ErrPasswordFormat = &Error{60,"password format invalid ","middle",nil}
)

// Error translate
type Error struct {
	Code   int32
	Detail string
	Prefix string
	vars   []interface{}
}

func (e *Error) GetPrefix() string {
	return e.Prefix
}

func (e *Error) GetCode() int32 {
	return e.Code
}

func (e *Error) GetKey() string {
	return fmt.Sprintf("prefix:%s , code:%d", e.GetPrefix(), e.GetCode())
}

func (e *Error) GetDetail() string {
	return fmt.Sprintf("detail: %s", e.Detail)
}

func (e *Error) GetVars() []interface{} {
	return e.vars
}

func (e *Error) SetDetail(c string) *Error {
	e.Detail = c
	return e
}

func (e *Error) SetVars(con ...interface{}) *Error {
	e2 := &Error{}
	*e2 = *e
	e2.vars = con
	return e2
}

func (e *Error) Error() string {
	return e.GetKey() + " , " + e.GetDetail()
}

// Text translate
type Text struct {
	Code   int32
	Detail string
	Prefix string
	vars   []interface{}
}

func (e *Text) GetPrefix() string {
	return e.Prefix
}

func (e *Text) GetCode() int32 {
	return e.Code
}

func (e *Text) GetKey() string {
	return fmt.Sprintf("%s%d", e.GetPrefix(), e.GetCode())
}

func (e *Text) GetDetail() string {
	return e.Detail
}

func (e *Text) GetVars() []interface{} {
	return e.vars
}

func (e *Text) SetDetail(c string) *Text {
	e.Detail = c
	return e
}

func (e *Text) SetVars(con ...interface{}) *Text {
	e2 := &Text{}
	*e2 = *e
	e2.vars = con
	return e2
}

func (e *Text) String() string {
	return fmt.Sprintf("%s%d|%s", e.Prefix, e.Code, e.Detail)
}