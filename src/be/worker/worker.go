package worker

import (
	"be/common"
	xe "be/common/error"
	"be/common/log"
	"be/options"
	"fmt"
	"path"
	"sync"

	uuid "github.com/satori/go.uuid"
)

type Worker struct {
	id        string
	serviceId string
	mgr       *WorkerMgr

	lock sync.Mutex

	// 源码路径
	codePath string
	// 项目全名
	projectFullName string

	// 是否正在下载源码
	inCodesFetching bool

	// 是否已经Open过了
	hasOpened bool
	// 字符串类型的配置信息
	rawConfig string
}

func NewWorker(serviceId string, mgr *WorkerMgr, codeDir string, fullName string) *Worker {
	return &Worker{
		id:              fmt.Sprintf("%s", uuid.NewV4()),
		serviceId:       serviceId,
		mgr:             mgr,
		codePath:        path.Join(options.Options.CodesRootPath, codeDir),
		projectFullName: fullName,
		lock:            sync.Mutex{},
		inCodesFetching: false,
		hasOpened:       false,
	}
}

// Die 删除此Worker
func (w *Worker) Die() {
	// todo
}

// FetchCodes 下载源码
func (w *Worker) FetchCodes() error {
	w.lock.Lock()
	if w.inCodesFetching {
		w.lock.Unlock()
		log.Errorln("已经处于源码下载状态")
		return xe.New("已经处于源码下载状态")
	}
	w.inCodesFetching = true
	w.lock.Unlock()

	// 执行下载操作
	codePath := w.codePath
	err := common.Mkdir(codePath)
	if err != nil {
		w.inCodesFetching = false
		return err
	}
	cmd := "git"
	args := []string{"clone", "--quiet", fmt.Sprintf("https://github.com/%s.git", w.projectFullName), codePath}

	o, e, err := common.Exec(options.Options.CodeFetchTimeout, cmd, args...)
	if err != nil {
		log.Errorf("下载源码失败: %s", err.Error())
		w.inCodesFetching = false
		return err
	}

	log.Debugf("下载源码stdout %s", o)
	log.Debugf("下载源码stderr %s", e)

	w.inCodesFetching = false
	return nil
}

// Open 打开项目
func (w *Worker) Open(rawConfig string) error {
	log.Debugf("worker %s 开始打开", w.id)
	w.lock.Lock()
	if w.hasOpened == true {
		w.lock.Unlock()
		return xe.New("worker已经处于Open状态")
	}
	w.hasOpened = true
	w.lock.Unlock()

	w.rawConfig = rawConfig

	// todo 启动各个分析服务

	return nil
}
