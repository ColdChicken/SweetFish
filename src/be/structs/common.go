package structs

import (
	"be/lex"
)

// UserInfo 用户信息，这个结构体用于管理控制台使用
type UserInfo struct {
	Username string `json:"username"`
}

// GitHubProject 一个GitHub项目的基本信息
type GitHubProject struct {
	Id             int64  `json:"id"`
	Name           string `json:"name"`
	FullName       string `json:"fullName"`
	HtmlUrl        string `json:"htmlUrl"`
	Description    string `json:"description"`
	CloneUrl       string `json:"cloneUrl"`
	HomePage       string `json:"homepage"`
	Size           int64  `json:"size"`
	Stars          int64  `json:"stars"`
	Watchers       int64  `json:"watchers"`
	Forks          int64  `json:"forks"`
	Language       string `json:"language"`
	OwnerLoginName string `json:"ownerLoginName"`
}

// GitHubProjectDetail 一个GitHub项目的详细信息
type GitHubProjectDetail struct {
}

// GitHubProjectDownloadInfo 下载GitHub到本地后的信息
type GitHubProjectDownloadInfo struct {
}

// Project 用户项目基本信息
type Project struct {
	Id       int64  `json:"id"`
	FullName string `json:"fullName"`
	Config   string `json:"config"`
	Status   string `json:"status"`
}

// ProjectInfo 内部初始化用的项目信息
type ProjectInfo struct {
	// 基本信息
	Id           int64
	Username     string
	FullName     string
	SourceCodeIp string
	Config       string
	Status       string
	LangTypesStr string
}

// 打开项目结果信息
type OpenProjectResult struct {
	// 打开结果: 成功/失败
	Result string `json:"result"`
	// 失败说明
	ErrMsg string `json:"errMsg"`
}

type ProjectUri struct {
	// 目录/文件
	Type string `json:"type"`
	// 目录或文件名
	Name string `json:"name"`
	// 目录的话下层相关信息
	SubDirs []*ProjectUri `json:"subDirs"`
}

// 项目的目录信息
type ProjectCatalog struct {
	Catalog []*ProjectUri `json:"catalog"`
}

// 打开文件返回信息
type OpenFileResult struct {
	Name       string              `json:"name"`
	RawContent string              `json:"rawContent"`
	Lines     [][]*lex.Token `json:"lines"`
}
