package structs

// UserInfo 用户信息，这个结构体用于管理控制台使用
type UserInfo struct {
	Username string `json:"username"`
}

// User 用户信息，这个结构体用于实际的用户
type User struct {
}

// GitHubProject 一个GitHub项目的基本信息
type GitHubProject struct {
}

// GitHubProjectDetail 一个GitHub项目的详细信息
type GitHubProjectDetail struct {
}

// GitHubProjectDownloadInfo 下载GitHub到本地后的信息
type GitHubProjectDownloadInfo struct {
}

// Project 用户项目基本信息
type Project struct {
}

// ProjectDetail 用户的项目详情
type ProjectDetail struct {
}

// ActionResult 用户采取动作后的返回信息
type ActionResult struct {
}
