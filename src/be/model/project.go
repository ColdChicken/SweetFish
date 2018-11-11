package model

import (
	"be/common/log"
	"be/service"
	"be/structs"
)

type ProjectMgr struct {
	serviceMgr *service.ServiceMgr
	projectMgr *service.ProjectMgr
}

var Project *ProjectMgr

func init() {
	serviceMgr := service.NewServiceMgr()
	projectMgr := service.NewProjectMgr(serviceMgr)
	Project = &ProjectMgr{
		serviceMgr: serviceMgr,
		projectMgr: projectMgr,
	}
}

// CreateProject 创建项目
// 此方法调用后会立刻返回
func (m *ProjectMgr) CreateProject(requestUser *structs.UserInfo, projectId string) error {
	log.Debugf("开始创建项目 %s - %s", requestUser.Username, projectId)
	// 获取目标主机
	targetIp, err := m.serviceMgr.ChooseTargetServiceMachine()
	if err != nil {
		log.Errorln(err.Error())
		return err
	}

	// 创建项目
	project, err := m.projectMgr.CreateProjectForUser(requestUser.Username, targetIp, projectId)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}

	// 开始初始化
	go project.Init()

	return nil
}

// ListProjects 列出被查询用户的项目基本信息列表
func (m *ProjectMgr) ListProjects(requestUser *structs.UserInfo) ([]*structs.Project, error) {
	return nil, nil
}

// GetProjectDetail 查询项目详情
func (m *ProjectMgr) GetProjectDetail(requestUser *structs.UserInfo, projectId string) (*structs.ProjectDetail, error) {
	return nil, nil
}

// OpenProject 用户打开某个已经创建完成的项目
// 打开项目意味着:
// * 对于LSP的服务，此动作会负责执行相关初始化动作
// * 加载此项目的项目信息、文件目录信息等返回给用户
// * 根据项目情况启动后台语义分析等服务
func (m *ProjectMgr) OpenProject(requestUser *structs.UserInfo, projectId string) (*structs.ProjectDetail, error) {
	return nil, nil
}

// DoActionInProject 用户在项目中进行某种操作
// action表示用户进行的操作，根据操作的不同actionRawInfo中会包含不同的信息
func (m *ProjectMgr) DoActionInProject(requestUser *structs.UserInfo, projectId string, action string, actionRawInfo string) (*structs.ActionResult, error) {
	return nil, nil
}

// CloseProject 用户关闭某个已经打开的项目
// 关闭项目会根据实际情况决定是否关闭后台服务
func (m *ProjectMgr) CloseProject(requestUser *structs.UserInfo, projectId string) (*structs.ProjectDetail, error) {
	return nil, nil
}

// DeleteProject 删除项目
func (m *ProjectMgr) DeleteProject(requestUser *structs.UserInfo, projectId string) error {
	return nil
}
