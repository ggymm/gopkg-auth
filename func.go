package gopkg_auth

import (
	"github.com/pkg/errors"
	"time"
)

func NotInit() bool {
	return auth == nil
}

// GetTokenName 获取 token 的名称
func GetTokenName() string {
	return auth.tokenName
}

// GetDefaultTimeout 获取 token 的默认过期时间
func GetDefaultTimeout() time.Duration {
	return auth.tokenTimeout
}

// Login 登录
func Login(id int64, config ...LoginConfig) (string, error) {
	if NotInit() {
		return "", errors.New(ErrAuthNotInit)
	}
	var cfg = LoginConfig{
		Device:  "web",
		Timeout: auth.tokenTimeout,
	}
	if len(config) > 0 {
		cfg = config[0]
	}
	return auth.Login(id, cfg)
}

func Logout(id int64, device ...string) error {
	if NotInit() {
		return errors.New(ErrAuthNotInit)
	}
	return auth.Logout(id, device...)
}

// Check 检查 token 是否有效
func Check(token string) (bool, error) {
	if NotInit() {
		return false, errors.New(ErrAuthNotInit)
	}
	return auth.CheckToken(token)
}

func GetSession(token string) (any, error) {
	if NotInit() {
		return nil, errors.New(ErrAuthNotInit)
	}
	return auth.GetSessionData(token)
}

func SaveSession(id int64, value any) error {
	if NotInit() {
		return errors.New(ErrAuthNotInit)
	}
	return auth.SaveSessionData(id, value)
}
