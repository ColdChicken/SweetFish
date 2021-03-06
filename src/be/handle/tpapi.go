package handle

import (
	"be/common"
	xe "be/common/error"
	"be/common/log"
	"be/model"
	"be/structs"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func tpapiGetRequestUser(req *http.Request) (*structs.UserInfo, error) {
	token := req.Header.Get("token")
	userinfo, err := model.Auth.GetUserInfoByToken(token)
	if err != nil {
		return nil, err
	}
	return userinfo, nil
}

func tpapiUserLogin(res http.ResponseWriter, req *http.Request) {
	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.WithFields(log.Fields{}).Error("请求报文解析失败")
		common.ResInvalidRequestBody(res)
		return
	}

	type Request struct {
		Code string `json:"code"`
	}

	request := &Request{}
	if err := common.ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("解析模板JSON失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}

	result, err := model.Auth.GetTokenByTPCode(request.Code)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("model请求处理失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}
	common.ResMsg(res, 200, result)

}

func tpapiSearchGithubProject(res http.ResponseWriter, req *http.Request) {
	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.WithFields(log.Fields{}).Error("请求报文解析失败")
		common.ResInvalidRequestBody(res)
		return
	}

	type Request struct {
		Inputs string `json:"inputs"`
	}

	request := &Request{}
	if err := common.ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("解析模板JSON失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}

	result, err := model.GitHub.SearchProjectsByFilter(request.Inputs)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("model请求处理失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}
	b, err := json.Marshal(result)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("JSON生成失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}
	common.ResMsg(res, 200, string(b))
}

func tpapiShowGitHubProjectDetail(res http.ResponseWriter, req *http.Request) {
	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.WithFields(log.Fields{}).Error("请求报文解析失败")
		common.ResInvalidRequestBody(res)
		return
	}

	type Request struct {
		ProjectId string `json:"projectId"`
	}

	request := &Request{}
	if err := common.ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("解析模板JSON失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}

	result, err := model.GitHub.ShowProjectDetail(request.ProjectId)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("model请求处理失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}
	b, err := json.Marshal(result)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("JSON生成失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}
	common.ResMsg(res, 200, string(b))
}

func tpapiCreateProject(res http.ResponseWriter, req *http.Request) {
	requestUser, err := tpapiGetRequestUser(req)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("获取用户信息失败")
		common.ResMsg(res, 400, xe.AuthError().Error())
		return
	}

	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.WithFields(log.Fields{}).Error("请求报文解析失败")
		common.ResInvalidRequestBody(res)
		return
	}

	type Request struct {
		ProjectId string `json:"projectId"`
	}

	request := &Request{}
	if err := common.ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("解析模板JSON失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}

	err = model.Project.CreateProject(requestUser, request.ProjectId)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("model请求处理失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}

	common.ResSuccessMsg(res, 200, "操作成功")
}

func tpapiListProjects(res http.ResponseWriter, req *http.Request) {
	requestUser, err := tpapiGetRequestUser(req)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("获取用户信息失败")
		common.ResMsg(res, 400, xe.AuthError().Error())
		return
	}

	result, err := model.Project.ListProjects(requestUser)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("model请求处理失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}
	b, err := json.Marshal(result)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("JSON生成失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}
	common.ResMsg(res, 200, string(b))
}

func tpapiListProjectCatalog(res http.ResponseWriter, req *http.Request) {
	requestUser, err := tpapiGetRequestUser(req)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("获取用户信息失败")
		common.ResMsg(res, 400, xe.AuthError().Error())
		return
	}

	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.WithFields(log.Fields{}).Error("请求报文解析失败")
		common.ResInvalidRequestBody(res)
		return
	}

	type Request struct {
		ProjectId int64 `json:"projectId"`
	}

	request := &Request{}
	if err := common.ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("解析模板JSON失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}

	result, err := model.Project.ListProjectCatalog(requestUser, request.ProjectId)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("model请求处理失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}
	b, err := json.Marshal(result)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("JSON生成失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}
	common.ResMsg(res, 200, string(b))
}

func tpapiActionOpenFile(res http.ResponseWriter, req *http.Request) {
	requestUser, err := tpapiGetRequestUser(req)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("获取用户信息失败")
		common.ResMsg(res, 400, xe.AuthError().Error())
		return
	}

	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.WithFields(log.Fields{}).Error("请求报文解析失败")
		common.ResInvalidRequestBody(res)
		return
	}

	type Request struct {
		ProjectId int64 `json:"projectId"`
		FilePath string `json:"filePath"`
		FileName string `json:"fileName"`
	}

	request := &Request{}
	if err := common.ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("解析模板JSON失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}

	result, err := model.Project.ActionOpenFile(requestUser, request.ProjectId, request.FilePath, request.FileName)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("model请求处理失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}
	b, err := json.Marshal(result)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("JSON生成失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}
	common.ResMsg(res, 200, string(b))
}

func tpapiOpenProject(res http.ResponseWriter, req *http.Request) {
	requestUser, err := tpapiGetRequestUser(req)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("获取用户信息失败")
		common.ResMsg(res, 400, xe.AuthError().Error())
		return
	}

	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.WithFields(log.Fields{}).Error("请求报文解析失败")
		common.ResInvalidRequestBody(res)
		return
	}

	type Request struct {
		ProjectId int64 `json:"projectId"`
	}

	request := &Request{}
	if err := common.ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("解析模板JSON失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}

	result, err := model.Project.OpenProject(requestUser, request.ProjectId)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("model请求处理失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}
	b, err := json.Marshal(result)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("JSON生成失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}
	common.ResMsg(res, 200, string(b))
}

func tpapiCloseProject(res http.ResponseWriter, req *http.Request) {
	requestUser, err := tpapiGetRequestUser(req)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("获取用户信息失败")
		common.ResMsg(res, 400, xe.AuthError().Error())
		return
	}

	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.WithFields(log.Fields{}).Error("请求报文解析失败")
		common.ResInvalidRequestBody(res)
		return
	}

	type Request struct {
		ProjectId int64 `json:"projectId"`
	}

	request := &Request{}
	if err := common.ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("解析模板JSON失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}

	model.Project.CloseProject(requestUser, request.ProjectId)
	common.ResSuccessMsg(res, 200, "操作成功")
}

func tpapiDeleteProject(res http.ResponseWriter, req *http.Request) {
	requestUser, err := tpapiGetRequestUser(req)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("获取用户信息失败")
		common.ResMsg(res, 400, xe.AuthError().Error())
		return
	}

	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.WithFields(log.Fields{}).Error("请求报文解析失败")
		common.ResInvalidRequestBody(res)
		return
	}

	type Request struct {
		ProjectId int64 `json:"projectId"`
	}

	request := &Request{}
	if err := common.ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("解析模板JSON失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}

	err = model.Project.DeleteProject(requestUser, request.ProjectId)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("model请求处理失败")
		common.ResMsg(res, 400, xe.HandleRequestError().Error())
		return
	}

	common.ResSuccessMsg(res, 200, "操作成功")
}
