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
	// 下载源码
	r.RegistURLMapping("/v1/worker/fetchcodes", "POST", apiFetchCodes)
	// 打开项目
	r.RegistURLMapping("/v1/worker/open", "POST", apiOpenProject)
}
