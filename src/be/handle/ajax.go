package handle

import (
	"be/common"
	xe "be/common/error"
	"be/common/log"
	"be/model"
	"be/session"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func tokenValidation(req *http.Request) error {
	token, err := session.CM.Get("token", req)
	if err != nil || token == "" {
		return xe.AuthError()
	}
	_, err = model.Auth.GetUserInfoByToken(token)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("根据token获取用户信息失败")
		return xe.AuthError()
	}
	return nil
}

func ajaxLogout(res http.ResponseWriter, req *http.Request) {
	session.CM.Remove("token", res)
	http.Redirect(res, req, "/ht", http.StatusTemporaryRedirect)
}

func ajaxGenTokenByUMAndPassword(res http.ResponseWriter, req *http.Request) {
	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.WithFields(log.Fields{}).Error("请求报文解析失败")
		common.ResInvalidRequestBody(res)
		return
	}

	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	request := &Request{}
	if err := common.ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("解析模板JSON失败")
		common.ResMsg(res, 400, err.Error())
		return
	}

	token, err := model.Auth.GenTokenByUsernameAndPassword(request.Username, request.Password)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("GenToken失败")
		common.ResMsg(res, 400, err.Error())
		return
	}

	// 在session中记录token
	session.CM.Set("token", token, res)
	common.ResSuccessMsg(res, 200, "token生成成功")
}

func ajaxGetUserInfo(res http.ResponseWriter, req *http.Request) {
	token, err := session.CM.Get("token", req)
	if err != nil || token == "" {
		common.ResMsg(res, 400, "请求中未包含token")
		return
	}
	userInfo, err := model.Auth.GetUserInfoByToken(token)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("根据token获取用户信息失败")
		common.ResMsg(res, 400, err.Error())
		return
	}
	b, err := json.Marshal(userInfo)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("根据token获取用户信息失败 JSON生成失败")
		common.ResMsg(res, 500, err.Error())
		return
	}
	common.ResMsg(res, 200, string(b))
}
