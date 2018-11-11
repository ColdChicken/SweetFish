package service

import (
	xe "be/common/error"
)

// 服务主机，Ip必须唯一
type Machine struct {
	Ip     string
	Port   string
	Status string
}

type ServiceMgr struct {
	// 目标服务主机
	machines []*Machine
}

func NewServiceMgr() *ServiceMgr {
	mgr := &ServiceMgr{
		machines: []*Machine{},
	}

	// 新增测试用主机
	mgr.machines = append(mgr.machines, &Machine{
		Ip:     "127.0.0.1",
		Port:   "19512",
		Status: "正常",
	})

	return mgr
}

// ChooseTargetServiceMachine 获取一台目标主机
func (m *ServiceMgr) ChooseTargetServiceMachine() (string, error) {
	if len(m.machines) == 0 {
		return "", xe.New("无可用主机")
	}
	// 当前简单点，就默认的一台主机
	return m.machines[0].Ip, nil
}

// CreateService 在目标主机上启动项目关联的分析服务
// 其会在目标主机上启动服务，并在本地保存存根
func (m *ServiceMgr) CreateService(targetMachineIp string, project *Project) (*Service, error) {
	return nil, nil
}

// GetProjectService 获取某个项目关联的服务
func (m *ServiceMgr) GetProjectService(project *Project) (*Service, error) {
	return nil, nil
}
