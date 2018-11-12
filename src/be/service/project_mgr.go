package service

import (
	xe "be/common/error"
	"be/common/log"
	"be/dao"
	"fmt"
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

	return projectMgr
}

// InitProjectsFromDB 从DB中初始化项目信息。此方法必须在使用此model前调用
func (m *ProjectMgr) InitProjectsFromDB() {
	// 加载project
	projectsInDB, err := m.projectDao.SyncAllProjectsForInitial()
	if err != nil {
		log.Errorf("从DB加载项目信息失败 %s", err.Error())
		panic(fmt.Sprintf("从DB加载项目信息失败 %s", err.Error()))
	}
	for _, projectInDB := range projectsInDB {
		if _, ok := m.projects[projectInDB.Username]; ok == false {
			m.projects[projectInDB.Username] = []*Project{}
		}

		// 状态不为正常的统一重置为失败，用户会的删除重建项目
		if projectInDB.Status != "正常" {
			projectInDB.Status = "失败"
		}

		m.projects[projectInDB.Username] = append(m.projects[projectInDB.Username], NewProject(projectInDB.Id, projectInDB.Username, projectInDB.FullName, projectInDB.SourceCodeIp, projectInDB.Config, projectInDB.Status, m, m.serviceMgr))
	}
}

func (m *ProjectMgr) GetProjectByUserAndProjectId(username string, projectId int64) (*Project, error) {
	m.projectsLock.Lock()
	if projects, ok := m.projects[username]; ok == true {
		for _, project := range projects {
			if project.Id == projectId {
				m.projectsLock.Unlock()
				return project, nil
			}
		}
		m.projectsLock.Unlock()
		return nil, xe.New("项目不存在")
	} else {
		m.projectsLock.Unlock()
		return nil, xe.New("项目不存在")
	}
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
