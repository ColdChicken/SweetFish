package worker

import (
	"be/server"
)

func New() *server.WWWMux {
	return server.New()
}

func InitHandle(r *server.WWWMux) {
	// 创建worker
	r.RegistURLMapping("/v1/worker/create", "POST", apiCreateWorker)
	// 销毁worker
	r.RegistURLMapping("/v1/worker/delete", "POST", apiDeleteWorker)
	// 下载源码
	r.RegistURLMapping("/v1/worker/fetchcodes", "POST", apiFetchCodes)
	// 分析源码信息
	r.RegistURLMapping("/v1/worker/init", "POST", apiInit)
	// 打开项目
	r.RegistURLMapping("/v1/worker/open", "POST", apiOpenProject)
}
