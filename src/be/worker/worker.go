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

	// 是否正在下载源码
	inCodesFetching bool
}

func NewWorker(serviceId string, mgr *WorkerMgr) *Worker {
	return &Worker{
		id:              fmt.Sprintf("%s", uuid.NewV4()),
		serviceId:       serviceId,
		mgr:             mgr,
		lock:            sync.Mutex{},
		inCodesFetching: false,
	}
}

// Die 删除此Worker
func (w *Worker) Die() {
	// todo
}

// FetchCodes 下载源码
func (w *Worker) FetchCodes(codeDir string, fullName string) error {
	w.lock.Lock()
	if w.inCodesFetching {
		w.lock.Unlock()
		log.Errorln("已经处于源码下载状态")
		return xe.New("已经处于源码下载状态")
	}
	w.inCodesFetching = true
	w.lock.Unlock()

	// 执行下载操作
	codePath := path.Join(options.Options.CodesRootPath, codeDir)
	err := common.Mkdir(codePath)
	if err != nil {
		w.inCodesFetching = false
		return err
	}
	cmd := "git"
	args := []string{"clone", "--quiet", fmt.Sprintf("https://github.com/%s.git", fullName), codePath}

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
