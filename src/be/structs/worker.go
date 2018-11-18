package structs

import (
	"be/common"
)

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
	ServiceId string            `json:"serviceId"`
	Config    string            `json:"config"`
	LangTypes []common.LangType `json:"langTypes"`
}

type WorkerDeleteWorkerRequest struct {
	ServiceId string `json:"serviceId"`
}

type WorkerInitRequest struct {
	ServiceId string `json:"serviceId"`
}

type WorkerInitResponse struct {
	LangTypes []common.LangType `json:"langTypes"`
}

type WorkerListCatalogRequest struct {
	ServiceId string `json:"serviceId"`
}

type WorkerListCatalogResponse struct {
	ProjectCatalog *ProjectCatalog `json:"catalog"`
}

type WorkerActionOpenFileRequest struct {
	ServiceId string `json:"serviceId"`
	FilePath string `json:"filePath"`
	FileName string `json:"fileName"`
}

type WorkerActionOpenFileResponse struct {
	OpenFileResult *OpenFileResult `json:"openFileResult"`
}
