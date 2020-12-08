package session

import (
	error2 "github.com/olongfen/contrib/utils"

	"net/http"
)

// 取出token字符串
func PubGetTokenFromReq(r *http.Request) (token string, err error) {
	if r == nil {
		err = error2.ErrSessionUndefined
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
					err = error2.ErrTokenReqHeaderOrFormKeyInvalid
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
		err = error2.ErrTokenChangePassword
		break
	case TokenCodeFreeze:
		err = error2.ErrTokenChangeFreeze
		break
	case TokenCodeLogout:
		err = error2.ErrTokenChangeLogout
		break
	default:
		err = error2.ErrTokenInvalid
		break
	}
	return
}
