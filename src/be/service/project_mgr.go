package service

import (
	"be/dao"
	"sync"
)

type ProjectMgr struct {
	projectDao *dao.ProjectDao
	serviceMgr *ServiceMgr

	// key: username, value: projects
	projectsLock sync.Mutex
	projects     map[string][]*Project
}

func NewProjectMgr(serviceMgr *ServiceMgr) *ProjectMgr {
	projectMgr := &ProjectMgr{
		projectDao:   &dao.ProjectDao{},
		serviceMgr:   serviceMgr,
		projectsLock: sync.Mutex{},
		projects:     map[string][]*Project{},
	}

	// todo 加载project

	return projectMgr
}

func (m *ProjectMgr) CreateProjectForUser(username string, sourceCodeIp string, fullName string) (*Project, error) {
	// 创建项目
	projectId, err := m.projectDao.CreateProject(fullName, sourceCodeIp)
	if err != nil {
		return nil, err
	}
	// 关联用户和项目
	err = m.projectDao.BindUserAndProject(projectId, username)
	if err != nil {
		return nil, err
	}
	// 获取项目信息
	id, fullName, status, sourceCodeIp, config, err := m.projectDao.GetProjectById(projectId)
	if err != nil {
		return nil, err
	}
	project := NewProject(id, username, fullName, sourceCodeIp, config, status, m, m.serviceMgr)
	m.projectsLock.Lock()
	if _, ok := m.projects[username]; ok == false {
		m.projects[username] = []*Project{}
	}
	m.projects[username] = append(m.projects[username], project)
	m.projectsLock.Unlock()
	return project, nil
}
