package model

import (
	xe "be/common/error"
	"be/common/log"
	"be/dao"
	"be/service"
	"be/structs"
)

type ProjectMgr struct {
	serviceMgr *service.ServiceMgr
	projectMgr *service.ProjectMgr
	projectDao *dao.ProjectDao
}

var Project *ProjectMgr

func init() {
	serviceMgr := service.NewServiceMgr()
	projectMgr := service.NewProjectMgr(serviceMgr)
	Project = &ProjectMgr{
		serviceMgr: serviceMgr,
		projectMgr: projectMgr,
		projectDao: &dao.ProjectDao{},
	}
}

func (m *ProjectMgr) InitProjects() {
	m.projectMgr.InitProjectsFromDB()
}

func (m *ProjectMgr) userQuotaCheck(requestUser *structs.UserInfo) (bool, error) {
	userProjects, err := m.ListProjects(requestUser)
	if err != nil {
		log.Errorln(err.Error())
		return false, err
	}
	// 目前一个用户只允许创建3个项目
	if len(userProjects) > 3 {
		return false, nil
	}
	return true, nil
}

// CreateProject 创建项目
// 此方法调用后会立刻返回
func (m *ProjectMgr) CreateProject(requestUser *structs.UserInfo, projectId string) error {
	log.Debugf("开始创建项目 %s - %s", requestUser.Username, projectId)
	// 判断用户是否还可以创建项目
	quotaCheck, err := m.userQuotaCheck(requestUser)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	if quotaCheck == false {
		return xe.New("当前用户添加的项目已超过上限")
	}

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
	return m.projectDao.ListProjectsByUser(requestUser.Username)
}

// OpenProject 用户打开某个已经创建完成的项目
// 打开项目意味着:
// * 对于LSP的服务，此动作会负责执行相关初始化动作
// * 根据项目情况启动后台语义分析等服务
// * ...
func (m *ProjectMgr) OpenProject(requestUser *structs.UserInfo, projectId int64) (*structs.OpenProjectResult, error) {
	project, err := m.projectMgr.GetProjectByUserAndProjectId(requestUser.Username, projectId)
	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}
	return project.Open()
}

// ListProjectCatalog 列出某个项目的目录
func (m *ProjectMgr) ListProjectCatalog(requestUser *structs.UserInfo, projectId int64) (*structs.ProjectCatalog, error) {
	project, err := m.projectMgr.GetProjectByUserAndProjectId(requestUser.Username, projectId)
	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}
	return project.ListCatalog()
}

// CloseProject 用户关闭某个已经打开的项目
// 关闭项目会根据实际情况决定是否关闭后台服务
func (m *ProjectMgr) CloseProject(requestUser *structs.UserInfo, projectId int64) {
	project, err := m.projectMgr.GetProjectByUserAndProjectId(requestUser.Username, projectId)
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	go project.Close()
}

// DeleteProject 删除项目
func (m *ProjectMgr) DeleteProject(requestUser *structs.UserInfo, projectId int64) error {
	project, err := m.projectMgr.GetProjectByUserAndProjectId(requestUser.Username, projectId)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}

	// 解除绑定关系
	m.projectMgr.UnbindUserAndProject(requestUser.Username, projectId)

	// 停止service
	go project.Remove()

	return nil
}

// ActionOpenFile 打开某个文件
func (m *ProjectMgr) ActionOpenFile(requestUser *structs.UserInfo, projectId int64, filePath string, fileName string) (*structs.OpenFileResult, error) {
	project, err := m.projectMgr.GetProjectByUserAndProjectId(requestUser.Username, projectId)
	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}
	return project.OpenFile(filePath, fileName)
}
