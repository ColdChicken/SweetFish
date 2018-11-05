package handle

import (
	"be/options"
	"be/server"
	"net/http"
)

/*
* InitHandle负责做Handle到实际URL的映射工作,
* handle包下的handle如果要被实际使用,则都需要在此进行注册
 */
func InitHandle(r *server.WWWMux) {
	// 初始化静态文件路径
	initStaticFileMapping(r)
	// 初始化管理控制台相关页面
	initAdminPortalMapping(r)
	// 初始化管理控制台ajax
	initAjaxMapping(r)
	// 小程序(tiny program)api相关的接口
	initTPAPIMapping(r)
}

func initStaticFileMapping(r *server.WWWMux) {
	fs := http.FileServer(http.Dir(options.Options.StaticFilePath))
	r.GetRouter().PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))
}

func initAdminPortalMapping(r *server.WWWMux) {
	// 后台
	r.RegistURLMapping("/ht", "GET", showIndexHtml)
	r.RegistURLMapping("/ht/ologin.html", "GET", showLoginHtml)
}

func initAjaxMapping(r *server.WWWMux) {
	// 注销
	r.RegistURLMapping("/v1/ajax/auth/logout", "GET", ajaxLogout)
	// 用户认证密码并生成token
	r.RegistURLMapping("/v1/ajax/auth/token", "POST", ajaxGenTokenByUMAndPassword)
	// 获取用户信息
	r.RegistURLMapping("/v1/ajax/auth/info", "GET", ajaxGetUserInfo)

}

func initTPAPIMapping(r *server.WWWMux) {
	// 根据用户输入搜索相关项目
	r.RegistURLMapping("/v1/tp/github/search", "POST", tpapiSearchGithubProject)
	// 根据用户选择展示项目详情
	r.RegistURLMapping("/v1/tp/github/projectdetail", "POST", tpapiShowGitHubProjectDetail)

	// 用户选择项目后创建相关项目
	r.RegistURLMapping("/v1/tp/project/create", "POST", tpapiCreateProject)
	// 获取项目相关信息
	r.RegistURLMapping("/v1/tp/project/detail", "POST", tpapiGetProjectDetail)
	// 列出当前用户的项目列表
	r.RegistURLMapping("/v1/tp/project/list", "POST", tpapiListProjects)
	// 用户打开项目
	r.RegistURLMapping("/v1/tp/project/open", "POST", tpapiOpenProject)
	// 用户在项目中执行某个操作
	r.RegistURLMapping("/v1/tp/project/action", "POST", tpapiDoActionInProject)
	// 用户关闭某个项目
	r.RegistURLMapping("/v1/tp/project/close", "POST", tpapiCloseProject)
	// 用户删除某个项目
	r.RegistURLMapping("/v1/tp/project/delete", "POST", tpapiDeleteProject)
}
