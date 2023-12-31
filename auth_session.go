package gopkg_auth

import (
	"github.com/ggymm/gopkg/utils"
)

func (a *Auth) GetSession(id int64, create bool) (*Session, error) {
	var (
		err       error
		sess      []byte
		session   *Session
		sessionId = a.sessionId(id)
	)

	sess, err = a.store.Get(sessionId)
	if err != nil {
		return nil, err
	}

	if sess == nil && create {
		return &Session{
			SessionId:  sessionId,
			UserId:     id,
			TokenList:  []Token{},
			CreateTime: utils.Now(),
		}, nil
	}

	utils.JsonDecodes(sess, &session)
	return session, nil
}

func (a *Auth) GetSessionData(token string) (any, error) {
	var (
		err     error
		value   []byte
		session *Session
	)

	// 获取 token 对应的 userId
	value, err = a.store.Get(a.tokenId(token))
	if err != nil || value == nil {
		return nil, err
	}
	userId, _, _ := a.parseTokenValue(value)

	// 获取 session
	session, err = a.GetSession(userId, false)
	if err != nil {
		return nil, err
	}

	return session.UserData, nil
}

func (a *Auth) SaveSessionData(id int64, value any) error {
	// 获取 session
	session, err := a.GetSession(id, false)
	if err != nil {
		return err
	}

	// 保存 session data
	return session.setUserData(value)
}

type Session struct {
	SessionId      []byte  // session id
	UserId         int64   // 用户 id
	UserData       any     // 用户自定义数据
	TokenList      []Token // token 列表
	CreateTime     int64   // 创建时间
	LastUpdateTime int64   // 最后更新时间
}

type Token struct {
	Value  string // token 值
	Device string // 设备信息
}

// 更新 session
func (s *Session) update() error {
	s.LastUpdateTime = utils.Now()
	return auth.store.Put(s.SessionId, utils.JsonEncode(s), NeverExpire)
}

func (s *Session) saveToken(token, device string) error {
	s.TokenList = append(s.TokenList, Token{
		Value:  token,
		Device: device,
	})
	return s.update()
}

func (s *Session) removeToken(token []string) error {
	if len(s.TokenList) == 0 {
		return nil
	}

	for _, t := range token {
		for i, v := range s.TokenList {
			if v.Value == t {
				s.TokenList = append(s.TokenList[:i], s.TokenList[i+1:]...)
			}
		}
	}
	return s.update()
}

func (s *Session) setUserData(value any) error {
	s.UserData = value
	return s.update()
}

func (s *Session) selectToken(device ...string) []string {
	var tokenList = make([]string, 0)
	if len(s.TokenList) == 0 {
		return tokenList
	}

	if len(device) > 0 {
		for _, t := range s.TokenList {
			for _, d := range device {
				if t.Device == d {
					tokenList = append(tokenList, t.Value)
				}
			}
		}
	} else {
		for _, t := range s.TokenList {
			tokenList = append(tokenList, t.Value)
		}
	}

	return tokenList
}
