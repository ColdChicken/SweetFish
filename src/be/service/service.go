package service

import (
	"be/common"
	xe "be/common/error"
	"be/common/log"
	"be/structs"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	uuid "github.com/satori/go.uuid"
)

type Service struct {
	id         string
	project    *Project
	serviceMgr *ServiceMgr

	// worker端主机
	workerMachine *Machine
	// 对端worker id
	workerId string
}

func NewService(project *Project, serviceMgr *ServiceMgr, machine *Machine) (*Service, error) {
	service := &Service{
		id:            fmt.Sprintf("%s", uuid.NewV4()),
		project:       project,
		serviceMgr:    serviceMgr,
		workerMachine: machine,
		workerId:      "",
	}

	return service, nil
}

// Connected 是否与对端取得链接
func (s *Service) Connected() bool {
	return s.workerId != ""
}

// ListCatalog 列出目录信息
func (s *Service) ListCatalog() (*structs.ProjectCatalog, error) {
	requestUrl := fmt.Sprintf("%s/v1/worker/listcatalog", s.workerMachine.getServiceAddress())

	request := &structs.WorkerListCatalogRequest{
		ServiceId: s.id,
	}

	hc := &http.Client{}

	b, err := json.Marshal(request)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("生成JSON串失败")
		return nil, err
	}

	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(b))
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("构造请求失败")
		return nil, err
	}

	resp, err := hc.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("发送请求失败")
		return nil, err
	}

	defer resp.Body.Close()
	respContent, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"msg": string(respContent),
		}).Error("请求对端处理失败")
		return nil, xe.New("请求对端处理失败")
	} else {
		response := &structs.WorkerListCatalogResponse{}
		if err := common.ParseJsonStr(string(respContent), response); err != nil {
			log.WithFields(log.Fields{
				"url":    requestUrl,
				"result": string(respContent),
				"err":    err.Error(),
			}).Error("解析模板JSON失败")
			return nil, xe.New("解析模板JSON失败")
		} else {
			return response.ProjectCatalog, nil
		}
	}
}

// Init 对端执行初始化动作
func (s *Service) Init() ([]common.LangType, error) {
	requestUrl := fmt.Sprintf("%s/v1/worker/init", s.workerMachine.getServiceAddress())

	request := &structs.WorkerInitRequest{
		ServiceId: s.id,
	}

	hc := &http.Client{}

	b, err := json.Marshal(request)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("生成JSON串失败")
		return nil, err
	}

	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(b))
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("构造请求失败")
		return nil, err
	}

	resp, err := hc.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("发送请求失败")
		return nil, err
	}

	defer resp.Body.Close()
	respContent, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"msg": string(respContent),
		}).Error("请求对端处理失败")
		return nil, xe.New("请求对端处理失败")
	} else {
		response := &structs.WorkerInitResponse{}
		if err := common.ParseJsonStr(string(respContent), response); err != nil {
			log.WithFields(log.Fields{
				"url":    requestUrl,
				"result": string(respContent),
				"err":    err.Error(),
			}).Error("解析模板JSON失败")
			return nil, xe.New("解析模板JSON失败")
		} else {
			return response.LangTypes, nil
		}
	}
}

// Connect 与对端创建worker服务
func (s *Service) Connect() {
	if s.workerId != "" {
		log.Warnln("当前service的workerId不为空，可能会造成未知问题")
	}

	requestUrl := fmt.Sprintf("%s/v1/worker/create", s.workerMachine.getServiceAddress())

	request := &structs.WorkerCreateWorkerRequest{ServiceId: s.id, CodeDir: s.project.GetCodeDir(), ProjectFullName: s.project.FullName}

	hc := &http.Client{}

	b, err := json.Marshal(request)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("生成JSON串失败")
		return
	}

	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(b))
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("构造请求失败")
		return
	}

	resp, err := hc.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("发送请求失败")
		return
	}

	defer resp.Body.Close()
	respContent, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"msg": string(respContent),
		}).Error("请求对端处理失败")
	} else {
		response := &structs.WorkerCreateWorkerResponse{}
		if err := common.ParseJsonStr(string(respContent), response); err != nil {
			log.WithFields(log.Fields{
				"url":    requestUrl,
				"result": string(respContent),
				"err":    err.Error(),
			}).Error("解析模板JSON失败")
		} else {
			s.workerId = response.WorkerId
		}
	}

}

// FetchCodes 下载源码
func (s *Service) FetchCodes() error {
	requestUrl := fmt.Sprintf("%s/v1/worker/fetchcodes", s.workerMachine.getServiceAddress())

	request := &structs.WorkerFetchCodesRequest{
		ServiceId: s.id,
	}

	hc := &http.Client{}

	b, err := json.Marshal(request)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("生成JSON串失败")
		return err
	}

	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(b))
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("构造请求失败")
		return err
	}

	resp, err := hc.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("发送请求失败")
		return err
	}

	defer resp.Body.Close()
	respContent, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"msg": string(respContent),
		}).Error("请求对端处理失败")
		return xe.New("请求对端处理失败")
	} else {
		return nil
	}
}

// Open 执行项目打开动作
func (s *Service) Open() error {
	requestUrl := fmt.Sprintf("%s/v1/worker/open", s.workerMachine.getServiceAddress())

	request := &structs.WorkerOpenProjectRequest{
		ServiceId: s.id,
		Config:    s.project.Config,
		LangTypes: s.project.langTypes,
	}

	hc := &http.Client{}

	b, err := json.Marshal(request)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("生成JSON串失败")
		return err
	}

	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(b))
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("构造请求失败")
		return err
	}

	resp, err := hc.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("发送请求失败")
		return err
	}

	defer resp.Body.Close()
	respContent, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"msg": string(respContent),
		}).Error("请求对端处理失败")
		return xe.New("请求对端处理失败")
	} else {
		return nil
	}
}

func (s *Service) RemoveDirs() error {
	// 删除对端的目录
	requestUrl := fmt.Sprintf("%s/v1/worker/removedirs", s.workerMachine.getServiceAddress())

	request := &structs.WorkerRemoveDirsRequest{
		ServiceId: s.id,
	}

	hc := &http.Client{}

	b, err := json.Marshal(request)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("生成JSON串失败")
		return err
	}

	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(b))
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("构造请求失败")
		return err
	}

	resp, err := hc.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("发送请求失败")
		return err
	}

	defer resp.Body.Close()
	respContent, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"msg": string(respContent),
		}).Error("请求对端处理失败")
		return xe.New("请求对端处理失败")
	} else {
		return nil
	}
}

// Remove 销毁此服务
func (s *Service) Remove() error {
	// 清除引用
	s.workerId = ""

	// 销毁关联的worker
	requestUrl := fmt.Sprintf("%s/v1/worker/delete", s.workerMachine.getServiceAddress())

	request := &structs.WorkerDeleteWorkerRequest{
		ServiceId: s.id,
	}

	hc := &http.Client{}

	b, err := json.Marshal(request)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("生成JSON串失败")
		return err
	}

	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(b))
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("构造请求失败")
		return err
	}

	resp, err := hc.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("发送请求失败")
		return err
	}

	defer resp.Body.Close()
	respContent, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"msg": string(respContent),
		}).Error("请求对端处理失败")
		return xe.New("请求对端处理失败")
	} else {
		return nil
	}
}

// OpenFile 打开文件
func (s *Service) OpenFile(filePath string, fileName string) (*structs.OpenFileResult, error) {
	requestUrl := fmt.Sprintf("%s/v1/worker/action/openfile", s.workerMachine.getServiceAddress())

	request := &structs.WorkerActionOpenFileRequest{
		ServiceId: s.id,
		FilePath:  filePath,
		FileName:  fileName,
	}

	hc := &http.Client{}

	b, err := json.Marshal(request)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("生成JSON串失败")
		return nil, err
	}

	req, err := http.NewRequest("POST", requestUrl, bytes.NewBuffer(b))
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("构造请求失败")
		return nil, err
	}

	resp, err := hc.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"err": err.Error(),
		}).Error("发送请求失败")
		return nil, err
	}

	defer resp.Body.Close()
	respContent, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		log.WithFields(log.Fields{
			"url": requestUrl,
			"msg": string(respContent),
		}).Error("请求对端处理失败")
		return nil, xe.New("请求对端处理失败")
	} else {
		response := &structs.WorkerActionOpenFileResponse{}
		if err := common.ParseJsonStr(string(respContent), response); err != nil {
			log.WithFields(log.Fields{
				"url":    requestUrl,
				"result": string(respContent),
				"err":    err.Error(),
			}).Error("解析模板JSON失败")
			return nil, xe.New("解析模板JSON失败")
		} else {
			return response.OpenFileResult, nil
		}
	}
}
