package service

import (
	"be/common"
	xe "be/common/error"
	"be/common/log"
	"be/dao"
	"be/structs"
	"encoding/json"
	"fmt"
	"path"
	"strings"
)

type Project struct {
	// 基本信息
	Id           int64
	Username     string
	FullName     string
	SourceCodeIp string
	Config       string
	Status       string

	// mgr的引用
	projectMgr *ProjectMgr
	serviceMgr *ServiceMgr

	// dao
	projectDao *dao.ProjectDao

	// 源码文件信息
	langTypes []common.LangType

	// 工作service
	service *Service
}

func NewProject(id int64, username string, fullName string, sourceCodeIp string, config string, status string, projectMgr *ProjectMgr, serviceMgr *ServiceMgr, langTypes []common.LangType) *Project {
	project := &Project{
		Id:           id,
		Username:     username,
		FullName:     fullName,
		SourceCodeIp: sourceCodeIp,
		Config:       config,
		Status:       status,
		projectMgr:   projectMgr,
		serviceMgr:   serviceMgr,
		projectDao:   &dao.ProjectDao{},
		langTypes:    langTypes,
	}
	// 默认的源码文件类型
	project.langTypes = append(project.langTypes, common.PlainText)
	return project
}

// GetCodeDir 获取源码的相对路径，对于同一个项目来说这个路径是唯一的
func (p *Project) GetCodeDir() string {
	return path.Join(fmt.Sprintf("%d", p.Id), strings.Replace(p.FullName, "/", "_", -1))
}

// Init 执行初始化动作
// 包含下载源码、执行初始化动作等
func (p *Project) Init() {
	service, err := p.serviceMgr.CreateService(p.SourceCodeIp, p)
	if err != nil {
		log.Errorf("project id %d, 获取service失败 %s", p.Id, err.Error())
		p.updateStatus("失败")
		return
	}
	p.setService(service)

	// 下载源码
	p.updateStatus("获取源码中")
	err = p.fetchCodes()
	if err != nil {
		log.Errorf("%d 获取下载源码失败 %s", p.Id, err.Error())
		p.updateStatus("失败")
		return
	}
	p.updateStatus("源码获取完成")

	// 执行初始化动作
	p.updateStatus("初始化中")
	err = p.init()
	if err != nil {
		log.Errorf("%d init失败 %s", p.Id, err.Error())
		p.updateStatus("失败")
		return
	}
	p.updateStatus("正常")
}

// fetchCodes 下载源码
func (p *Project) fetchCodes() error {
	err := p.service.FetchCodes()
	if err != nil {
		log.Errorln(err.Error())
	}
	return err
}

// init 初始化项目
func (p *Project) init() error {
	// 获取此项目的源码文件信息
	langTypes, err := p.service.Init()
	if err != nil {
		return err
	}
	p.langTypes = langTypes

	// 相关信息入库
	err = p.syncLangTypes()
	if err != nil {
		log.Errorf(err.Error())
		return err
	}

	return nil
}

func (p *Project) syncLangTypes() error {
	b, err := json.Marshal(p.langTypes)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("JSON生成失败")
		return err
	}
	return p.projectDao.SyncLangTypes(p.Id, string(b))
}

func (p *Project) updateStatus(status string) error {
	p.Status = status
	err := p.projectDao.UpdateStatus(p.Id, status)
	if err != nil {
		log.Errorf(err.Error())
	}
	return err
}

func (p *Project) setService(service *Service) {
	if p.service != nil {
		p.service.Remove()
		p.service = nil
	}
	p.service = service
}

// ListCatalog 列出项目目录
func (p *Project) ListCatalog() (*structs.ProjectCatalog, error) {
	if p.Status != "正常" {
		return nil, xe.New("项目状态不处于 正常 状态")
	}

	if p.service == nil {
		return nil, xe.New("service不存在")
	}

	return p.service.ListCatalog()
}

// Close 关闭项目
func (p *Project) Close() {
	if p.Status != "正常" {
		log.Errorf("%d 的状态不为 正常，不能关闭", p.Id)
	}
	if p.service != nil {
		p.service.Remove()
	}
}

// Remove 删除项目
// todo 这里删除项目的逻辑存在问题，代码不会被真正删除
func (p *Project) Remove() {
	if p.service != nil {
		p.service.Remove()
	}
}

// Open 打开项目
func (p *Project) Open() (*structs.OpenProjectResult, error) {
	result := &structs.OpenProjectResult{Result: "成功"}

	// 判断项目状态是否可以打开
	if p.Status != "正常" {
		result.Result = "失败"
		result.ErrMsg = "当前的项目状态不为'正常'。如果项目正处于'加载中'则请等待服务器端下载完成相关源码，否则请删除项目后重新添加此项目。"
		return result, nil
	}

	// 使用新的service
	if p.service != nil {
		p.service.Remove()
		p.service = nil
	}
	service, err := p.serviceMgr.CreateService(p.SourceCodeIp, p)
	if err != nil {
		log.Errorf("project id %d, 获取service失败 %s", p.Id, err.Error())
		p.updateStatus("失败")
		result.Result = "失败"
		result.ErrMsg = "服务端异常，请稍后重试"
		return result, nil
	}
	p.setService(service)

	// worker端执行open动作
	p.updateStatus("打开中")
	err = p.service.Open()
	if err != nil {
		log.Errorf("%d Open失败 %s", p.Id, err.Error())
		p.updateStatus("正常")
		result.Result = "失败"
		result.ErrMsg = "服务端异常，请稍后重试"
		return result, nil
	}
	p.updateStatus("正常")
	result.Result = "成功"
	return result, nil
}

func (p *Project) OpenFile(filePath string, fileName string) (*structs.OpenFileResult, error) {
	if p.Status != "正常" {
		return nil, xe.New("项目状态不处于 正常 状态")
	}

	if p.service == nil {
		return nil, xe.New("service不存在")
	}

	return p.service.OpenFile(filePath, fileName)
}

func (p *Project) RemoveDirs() error {
	if p.service != nil {
		return xe.New("RemoveDirs 只能一次性使用")
	}
	service, err := p.serviceMgr.CreateService(p.SourceCodeIp, p)
	if err != nil {
		return err
	}

	err = service.RemoveDirs()
	if err != nil {
		return err
	}

	service.Remove()

	p.updateStatus("已删除")
	return nil
}
