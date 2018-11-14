package process

type ProcessMgr struct {
}

func NewProcessMgr() *ProcessMgr {
	m := &ProcessMgr{}
	return m
}

// Close 关闭其所管理的进程
func (m *ProcessMgr) Close() {

}
