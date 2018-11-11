package structs

type WorkerErrMsg struct {
	Msg string `json:"msg"`
}

type WorkerCreateWorkerRequest struct {
	ServiceId string `json:"serviceId"`
}

type WorkerCreateWorkerResponse struct {
	WorkerId string `json:"workerId"`
}

type WorkerFetchCodesRequest struct {
	ServiceId       string `json:"serviceId"`
	CodeDir         string `json:"codeDir"`
	ProjectFullName string `json:"projectFullName"`
}
