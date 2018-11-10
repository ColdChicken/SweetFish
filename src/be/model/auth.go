package model

import (
	"be/common"
	xe "be/common/error"
	"be/common/log"
	"be/dao"
	"be/options"
	"be/structs"
	"fmt"
	"io/ioutil"
	"net/http"
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

func (m *AuthMgr) GetTokenByTPCode(code string) (string, error) {
	requestUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", options.Options.TPAppId, options.Options.TPSecret, code)

	hc := &http.Client{}

	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("构造请求失败")
		return "", err
	}

	resp, err := hc.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("构造请求失败")
		return "", err
	}

	defer resp.Body.Close()
	respContent, _ := ioutil.ReadAll(resp.Body)

	type wxResponse struct {
		OpenId  string `json:"openid"`
		ErrCode int64  `json:"errcode"`
		ErrMsg  string `json:"errMsg"`
	}

	response := &wxResponse{}
	if err := common.ParseJsonStr(string(respContent), response); err != nil {
		log.WithFields(log.Fields{
			"result": string(respContent),
			"err":    err.Error(),
		}).Error("解析模板JSON失败")
		return "", err
	}

	if response.ErrCode != 0 {
		log.Errorf("获取openid失败， err: %s", response.ErrMsg)
		return "", xe.New(response.ErrMsg)
	}

	// 生成用户信息
	userId, err := m.dao.CreateUserByOpenId(response.OpenId)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("创建用户信息失败")
		return "", err
	}

	// 获取token
	token, err := m.dao.GenTokenByUsernameAndPassword(userId, response.OpenId)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("获取用户token失败")
		return "", err
	}

	return token, nil
}
