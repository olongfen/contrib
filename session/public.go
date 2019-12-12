package session

import (
	project "github.com/srlemon/contrib"

	"net/http"
)

// 取出token字符串
func PubGetTokenFromReq(r *http.Request) (token string, err error) {
	if r == nil {
		err = project.ErrSessionUndefined
		return
	}

	// format of token in header like "Bearer xxx"
	if token = r.Header.Get("Authentication"); len(token) > 7 {
		token = token[7:]
	} else {
		if token = r.Header.Get("Authorization"); len(token) > 7 {
			token = token[7:]
		}
	}

	// 兼容

	// url cover header define
	if v := r.FormValue("_token"); len(v) > 0 {
		token = v
	}
	return
}

// token-code对应的错误
func PubTokenCodeError(code int8) (err error) {
	switch code {
	case TokenCodePass:
		break
	case TokenCodePsw:
		err = project.ErrTokenChangePassword
		break
	case TokenCodeFreeze:
		err = project.ErrTokenChangeFreeze
		break
	case TokenCodeLogout:
		err = project.ErrTokenChangeLogout
		break
	default:
		err = project.ErrTokenInvalid
		break
	}
	return
}
