package service

import (
	xe "be/common/error"
	"be/common/log"
	"fmt"
)

// 服务主机，Ip必须唯一
type Machine struct {
	Ip     string
	Port   string
	Status string
}

func (m *Machine) getServiceAddress() string {
	return fmt.Sprintf("http://%s:%s", m.Ip, m.Port)
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
	// todo 改为可配置
	mgr.machines = append(mgr.machines, &Machine{
		Ip:     "127.0.0.1",
		Port:   "19512",
		Status: "正常",
	})

	return mgr
}

func (m *ServiceMgr) getMachineByIp(ip string) *Machine {
	for _, machine := range m.machines {
		if machine.Ip == ip {
			return machine
		}
	}
	return nil
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
	machine := m.getMachineByIp(targetMachineIp)
	if machine == nil {
		log.Errorln("无法获取主机")
		return nil, xe.New("无法获取主机")
	}

	service, err := NewService(project, project.serviceMgr, machine)
	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}

	service.Connect()
	if service.Connected() {
		return service, nil
	} else {
		return nil, xe.New("无法启动对端worker")
	}
}

// GetProjectService 获取某个项目关联的服务
func (m *ServiceMgr) GetProjectService(project *Project) (*Service, error) {
	return nil, nil
}
