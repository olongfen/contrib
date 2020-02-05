package session

import (
	project "github.com/olongfen/contrib"

	"net/http"
)

// 取出token字符串
func PubGetTokenFromReq(r *http.Request) (token string, err error) {
	if r == nil {
		err = project.ErrSessionUndefined
		return
	}

	// format of token in header like "Bearer xxx"
	if token = r.Header.Get("Authentication"); len(token) > 0 {
		return
	} else {
		if token = r.Header.Get("Authorization"); len(token) > 0 {
			return
		} else {
			token = r.Header.Get("token")
			if len(token) > 0 {
				return
			} else {
				if v := r.FormValue("token"); len(v) > 0 {
					token = v
					return
				} else {
					err = project.ErrTokenReqHeaderOrFormKeyInvalid
					return
				}

				return
			}
		}
	}

	// 兼容

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
