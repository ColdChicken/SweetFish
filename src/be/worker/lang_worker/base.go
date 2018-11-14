package lang_worker

import (
	"be/worker/process"
)

type File struct {
}

type FilePath struct {
}

type Position struct {
}

type LangWorker interface {
	// 初始化
	Init(pm *process.ProcessMgr, codePath string, projectFullName string, rawConfig string)

	// 动作: 打开文件，这里File可以包含语法高亮信息
	OpenFile(p *FilePath) *File

	// 动作: 跳转到定义
	GoToDefinition(p *Position) []*Position

	// 动作: 跳转到类型定义
	GoToTypeDefition(p *Position) []*Position

	// 动作: 跳转到实现
	GoToImplementation(p *Position) []*Position

	// 销毁
	Close()
}
