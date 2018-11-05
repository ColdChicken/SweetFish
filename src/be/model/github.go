package model

import (
	"be/structs"
)

type GitHubMgr struct{}

var GitHub *GitHubMgr

func init() {
	GitHub = &GitHubMgr{}
}

// SearchProjectsByFilter 根据用户的搜索条件列出对应的GitHub项目
func (m *GitHubMgr) SearchProjectsByFilter(filter string) ([]*structs.GitHubProject, error) {
	return nil, nil
}

// ShowProjectDetail 获取一个项目的详细信息
func (m *GitHubMgr) ShowProjectDetail(projectId string) (*structs.GitHubProjectDetail, error) {
	return nil, nil
}

// DownloadProject 下载项目到本地
func (m *GitHubMgr) DownloadProject(projectId string) (*structs.GitHubProjectDownloadInfo, error) {
	return nil, nil
}
