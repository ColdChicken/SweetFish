package structs

type WorkerErrMsg struct {
	Msg string `json:"msg"`
}

type WorkerCreateWorkerRequest struct {
	ServiceId       string `json:"serviceId"`
	CodeDir         string `json:"codeDir"`
	ProjectFullName string `json:"projectFullName"`
}

type WorkerCreateWorkerResponse struct {
	WorkerId string `json:"workerId"`
}

type WorkerFetchCodesRequest struct {
	ServiceId string `json:"serviceId"`
}

type WorkerOpenProjectRequest struct {
	ServiceId string `json:"serviceId"`
	Config    string `json:"config"`
}
