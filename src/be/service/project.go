package service

import (
	"be/common/log"
	"be/dao"
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

	// 工作service
	service *Service
}

func NewProject(id int64, username string, fullName string, sourceCodeIp string, config string, status string, projectMgr *ProjectMgr, serviceMgr *ServiceMgr) *Project {
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
	}
	return project
}

// Init 执行初始化动作
// 包含下载源码、执行初始化动作等
func (p *Project) Init() {
	service, err := p.serviceMgr.CreateService(p.SourceCodeIp, p)
	if err != nil {
		log.Errorf("%d 获取service失败 %s", p.Id, err.Error())
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
	return nil
}

// init 初始化项目
func (p *Project) init() error {
	return nil
}

func (p *Project) updateStatus(status string) error {
	err := p.projectDao.UpdateStatus(p.Id, status)
	if err != nil {
		log.Errorf(err.Error())
	}
	return err
}

func (p *Project) setService(service *Service) {
	if p.service != nil {
		p.service.Remove()
	}
	p.service = service
}
