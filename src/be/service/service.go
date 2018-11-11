package service

type Service struct {
	project    *Project
	serviceMgr *ServiceMgr
}

func NewService(project *Project, serviceMgr *ServiceMgr) (*Service, error) {
	service := &Service{
		project:    project,
		serviceMgr: serviceMgr,
	}

	return service, nil
}

// Remove 销毁此服务
func (s *Service) Remove() {

}
