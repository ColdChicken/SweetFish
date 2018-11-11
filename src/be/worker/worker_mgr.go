package worker

import (
	xe "be/common/error"
	"sync"
)

type WorkerMgr struct {
	workersLock sync.Mutex
	workers     map[string]*Worker
}

var workerMgr *WorkerMgr

func init() {
	workerMgr = &WorkerMgr{
		workersLock: sync.Mutex{},
		workers:     map[string]*Worker{},
	}
}

// CreateWorker 创建worker
func (m *WorkerMgr) CreateWorker(serviceId string) (*Worker, error) {
	worker := NewWorker(serviceId, m)
	m.workersLock.Lock()
	if _, ok := m.workers[serviceId]; ok == false {
		m.workers[serviceId] = worker
	} else {
		m.workers[serviceId].Die()
		m.workers[serviceId] = worker
	}
	m.workersLock.Unlock()
	return worker, nil
}

// GetWorkerByServiceId 根据service获取绑定的worker
func (m *WorkerMgr) GetWorkerByServiceId(serviceId string) (*Worker, error) {
	m.workersLock.Lock()
	if w, ok := m.workers[serviceId]; ok == true {
		m.workersLock.Unlock()
		return w, nil
	} else {
		m.workersLock.Unlock()
		return nil, xe.New("worker不存在")
	}
}
