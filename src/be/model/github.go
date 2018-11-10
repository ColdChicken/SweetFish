package model

import (
	"be/common"
	"be/common/log"
	"be/structs"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type GitHubMgr struct{}

var GitHub *GitHubMgr

func init() {
	GitHub = &GitHubMgr{}
}

// SearchProjectsByFilter 根据用户的搜索条件列出对应的GitHub项目
func (m *GitHubMgr) SearchProjectsByFilter(filter string) ([]*structs.GitHubProject, error) {
	const searchApi = "https://api.github.com/search/repositories?q="
	searchUrl := fmt.Sprintf("%s%s", searchApi, strings.TrimSpace(filter))

	hc := &http.Client{}

	req, err := http.NewRequest("GET", searchUrl, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"url": searchUrl,
			"err": err.Error(),
		}).Error("构造请求失败")
		return nil, err
	}

	resp, err := hc.Do(req)
	if err != nil {
		log.WithFields(log.Fields{
			"url": searchUrl,
			"err": err.Error(),
		}).Error("构造请求失败")
		return nil, err
	}

	defer resp.Body.Close()
	respContent, _ := ioutil.ReadAll(resp.Body)

	// GitHub结果
	type githubSearchRecordOwner struct {
		Login string `json:"login"`
	}

	type githubSearchRecord struct {
		Id          int64                    `json:"id"`
		Name        string                   `json:"name"`
		FullName    string                   `json:"full_name"`
		HtmlUrl     string                   `json:"html_url"`
		Description string                   `json:"description"`
		CloneUrl    string                   `json:"clone_url"`
		HomePage    string                   `json:"homepage"`
		Size        int64                    `json:"size"`
		Stars       int64                    `json:"stargazers_count"`
		Watchers    int64                    `json:"watchers_count"`
		Forks       int64                    `json:"forks_count"`
		Language    string                   `json:"language"`
		Owner       *githubSearchRecordOwner `json:"owner"`
	}

	type githubSearchResult struct {
		Items []*githubSearchRecord `json:"items"`
	}

	result := &githubSearchResult{}
	if err := common.ParseJsonStr(string(respContent), result); err != nil {
		log.WithFields(log.Fields{
			"url":    searchUrl,
			"result": string(respContent),
			"err":    err.Error(),
		}).Error("解析模板JSON失败")
		return nil, err
	}

	projects := []*structs.GitHubProject{}
	for _, record := range result.Items {
		projects = append(projects, &structs.GitHubProject{
			Id:             record.Id,
			Name:           record.Name,
			FullName:       record.FullName,
			HtmlUrl:        record.HtmlUrl,
			Description:    record.Description,
			CloneUrl:       record.CloneUrl,
			HomePage:       record.HomePage,
			Size:           record.Size,
			Stars:          record.Stars,
			Watchers:       record.Watchers,
			Forks:          record.Forks,
			Language:       record.Language,
			OwnerLoginName: record.Owner.Login,
		})
	}

	return projects, nil
}

// ShowProjectDetail 获取一个项目的详细信息
func (m *GitHubMgr) ShowProjectDetail(projectId string) (*structs.GitHubProjectDetail, error) {
	return nil, nil
}

// DownloadProject 下载项目到本地
func (m *GitHubMgr) DownloadProject(projectId string) (*structs.GitHubProjectDownloadInfo, error) {
	return nil, nil
}
