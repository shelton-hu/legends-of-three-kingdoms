package gerr

import (
	"fmt"
)

type ErrorMessage struct {
	Name string
	en   string
	zhCN string
}

func (em ErrorMessage) Message(locale string, err error, a ...interface{}) string {
	format := ""
	switch locale {
	case En:
		format = em.en
	case ZhCN:
		if len(em.zhCN) > 0 {
			format = em.zhCN
		} else {
			format = em.en
		}
	}
	if err != nil {
		return fmt.Sprintf("%s: %s", fmt.Sprintf(format, a...), err.Error())
	} else {
		return fmt.Sprintf(format, a...)
	}
}

// IAMManager
var (
	ErrorPermissionDenied = ErrorMessage{
		Name: "permission_denied",
		en:   "permission denied",
		zhCN: "没有权限",
	}
	ErrorAuthFailure = ErrorMessage{
		Name: "auth_failure",
		en:   "auth failure",
		zhCN: "认证失败",
	}
	ErrorAccessTokenExpired = ErrorMessage{
		Name: "access_token_expired",
		en:   "access token expired",
		zhCN: "访问令牌已过期",
	}
	ErrorRefreshTokenExpired = ErrorMessage{
		Name: "refresh_token_expired",
		en:   "refresh token expired",
		zhCN: "刷新令牌已过期",
	}
	ErrorInternalError = ErrorMessage{
		Name: "internal_error",
		en:   "internal error",
		zhCN: "内部错误",
	}
	ErrorMissingParameter = ErrorMessage{
		Name: "missing_parameter",
		en:   "missing parameter [%s]",
		zhCN: "缺少参数[%s]",
	}
	ErrorUnsupportedParameterValue = ErrorMessage{
		Name: "unsupported_parameter_value",
		en:   "unsupported parameter [%s] value [%s]",
		zhCN: "参数[%s]不支持值[%s]",
	}
	ErrorParameterShouldNotBeEmpty = ErrorMessage{
		Name: "parameter_should_not_be_empty",
		en:   "parameter [%s] should not be empty",
		zhCN: "参数[%s]不应该为空",
	}
)
