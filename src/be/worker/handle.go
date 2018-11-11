package worker

import (
	"be/common"
	"be/common/log"
	"be/structs"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func apiCreateWorker(res http.ResponseWriter, req *http.Request) {
	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.WithFields(log.Fields{}).Error("请求报文解析失败")
		common.ResInvalidRequestBody(res)
		return
	}

	request := &structs.WorkerCreateWorkerRequest{}
	if err := common.ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("解析模板JSON失败")
		common.ResMsg(res, 400, err.Error())
		return
	}

	worker, err := workerMgr.CreateWorker(request.ServiceId)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("model处理请求失败")
		common.ResMsg(res, 400, err.Error())
		return
	}

	response := &structs.WorkerCreateWorkerResponse{
		WorkerId: worker.id,
	}

	b, err := json.Marshal(response)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("JSON生成失败")
		common.ResMsg(res, 500, err.Error())
		return
	}
	common.ResMsg(res, 200, string(b))
}

func apiFetchCodes(res http.ResponseWriter, req *http.Request) {
	reqContent, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		log.WithFields(log.Fields{}).Error("请求报文解析失败")
		common.ResInvalidRequestBody(res)
		return
	}

	request := &structs.WorkerFetchCodesRequest{}
	if err := common.ParseJsonStr(string(reqContent), request); err != nil {
		log.Errorln("解析模板JSON失败")
		common.ResMsg(res, 400, err.Error())
		return
	}

	worker, err := workerMgr.GetWorkerByServiceId(request.ServiceId)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("model处理请求失败")
		common.ResMsg(res, 400, err.Error())
		return
	}

	err = worker.FetchCodes(request.CodeDir, request.ProjectFullName)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("worker处理请求失败")
		common.ResMsg(res, 400, err.Error())
		return
	}
	common.ResSuccessMsg(res, 200, "操作成功")
}
