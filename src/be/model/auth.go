package model

import (
	"be/common/log"
	"be/dao"
	"be/structs"
)

type AuthMgr struct {
	dao *dao.AuthDao
}

var Auth *AuthMgr

func init() {
	Auth = &AuthMgr{
		dao: &dao.AuthDao{},
	}
}

func (m *AuthMgr) GenTokenByUsernameAndPassword(username string, password string) (string, error) {
	return m.dao.GenTokenByUsernameAndPassword(username, password)
}

func (m *AuthMgr) GetUserInfoByToken(token string) (*structs.UserInfo, error) {
	userInfo, err := m.dao.GetUserInfoByToken(token)
	if err != nil {
		log.WithFields(log.Fields{
			"token": token,
			"err":   err.Error(),
		}).Error("根据token获取用户信息失败")
		return nil, err
	}
	return userInfo, nil
}
