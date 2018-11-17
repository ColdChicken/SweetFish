package worker

import (
	"be/common"
	xe "be/common/error"
	"be/common/log"
	"be/options"
	"be/structs"
	"be/worker/lang_worker"
	"be/worker/process"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
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

	// language worker, key为语言名，如go、javascript
	langWorkers map[common.LangType]lang_worker.LangWorker

	// 进程管理器
	processMgr *process.ProcessMgr

	// 是否正在下载源码
	inCodesFetching bool

	// 是否已经Open过了
	hasOpened bool
	// 字符串类型的配置信息
	rawConfig string
	// key-value类型的配置信息
	kvConfig map[common.LangType]string
}

func NewWorker(serviceId string, mgr *WorkerMgr, codeDir string, fullName string) *Worker {
	worker := &Worker{
		id:              fmt.Sprintf("%s", uuid.NewV4()),
		serviceId:       serviceId,
		mgr:             mgr,
		codePath:        path.Join(options.Options.CodesRootPath, codeDir),
		projectFullName: fullName,
		lock:            sync.Mutex{},
		inCodesFetching: false,
		hasOpened:       false,
		processMgr:      process.NewProcessMgr(),
		kvConfig:        map[common.LangType]string{},
		langWorkers:     map[common.LangType]lang_worker.LangWorker{},
	}

	return worker
}

// Die 删除此Worker
// (todo) 此方法只会停止该Worker的相关资源
func (w *Worker) Die() {
	go w.die()
}
func (w *Worker) die() {
	// 调用各个langWorker的关闭动作
	for _, lw := range w.langWorkers {
		lw.Close()
	}

	// 关闭进程
	w.processMgr.Close()

	// 移除此worker
	w.mgr.RemoveWorker(w)
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
func (w *Worker) Open(rawConfig string, langTypes []common.LangType) error {
	log.Debugf("worker %s 开始打开", w.id)
	w.lock.Lock()
	if w.hasOpened == true {
		w.lock.Unlock()
		return xe.New("worker已经处于Open状态")
	}
	w.hasOpened = true
	w.lock.Unlock()

	w.rawConfig = rawConfig
	if strings.TrimSpace(w.rawConfig) == "" {
		w.rawConfig = "{}"
	}
	// 解析配置
	if err := common.ParseJsonStr(w.rawConfig, &w.kvConfig); err != nil {
		log.Errorln("配置格式不正确, 解析模板JSON失败")
		return xe.New("配置格式不正确")
	}

	// 获取语言服务器
	for _, langType := range langTypes {
		lw := lang_worker.GetLangWorkerByLangType(langType)
		if lw == nil {
			log.Errorln("无法获取到语言服务器")
			continue
		}
		w.langWorkers[langType] = lw
	}

	// 调用各个langWorker的启动初始化动作
	for lang, lw := range w.langWorkers {
		if cfg, ok := w.kvConfig[lang]; ok == true {
			lw.Init(w.processMgr, w.codePath, w.projectFullName, cfg)
		} else {
			lw.Init(w.processMgr, w.codePath, w.projectFullName, "")
		}

	}

	return nil
}

// Open 初始化
func (w *Worker) Init() (*structs.WorkerInitResponse, error) {
	response := &structs.WorkerInitResponse{LangTypes: []common.LangType{}}
	// 获取此项目的源码文件信息
	files, err := common.GetDirAllFiles(w.codePath)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}

	// 判断文件类型
	for _, fileName := range files {
		langType := common.GetLangTypeByFileName(fileName)
		hasLangType := false
		for _, lt := range response.LangTypes {
			if lt == langType {
				hasLangType = true
				break
			}
		}
		if hasLangType == false {
			response.LangTypes = append(response.LangTypes, langType)
		}
	}

	return response, nil
}

func (w *Worker) Inited() bool {
	w.lock.Lock()
	if w.inCodesFetching {
		w.lock.Unlock()
		return false
	}
	if w.hasOpened == false {
		w.lock.Unlock()
		return false
	}
	w.lock.Unlock()
	return true
}

func getDirCatalogInfo(dir string) (*structs.ProjectCatalog, error) {
	catalog := &structs.ProjectCatalog{Catalog: []*structs.ProjectUri{}}

	directory, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	// 先加载目录
	for _, f := range directory {
		if f.IsDir() {
			if f.Name() == ".git" {
				continue
			}
			uri := &structs.ProjectUri{
				Type:    "目录",
				Name:    f.Name(),
				SubDirs: []*structs.ProjectUri{},
			}
			subDirs, err := getDirCatalogInfo(path.Join(dir, f.Name()))
			if err != nil {
				return nil, err
			}
			uri.SubDirs = subDirs.Catalog
			catalog.Catalog = append(catalog.Catalog, uri)
		} else {
			continue
		}
	}
	// 后加载文件
	for _, f := range directory {
		if f.IsDir() {
			continue
		} else {
			uri := &structs.ProjectUri{
				Type:    "文件",
				Name:    f.Name(),
				SubDirs: []*structs.ProjectUri{},
			}
			catalog.Catalog = append(catalog.Catalog, uri)
		}
	}

	return catalog, nil
}

// ListCatalog 列出目录
func (w *Worker) ListCatalog() (*structs.WorkerListCatalogResponse, error) {
	response := &structs.WorkerListCatalogResponse{}
	// 判断状态
	if w.Inited() == false {
		return nil, xe.New("worker当前状态无法处理该请求")
	}
	// 获取目录信息
	catalog, err := getDirCatalogInfo(w.codePath)
	if err != nil {
		log.Errorf("获取目录信息失败 %s", err.Error())
		return nil, err
	}
	// 做个排序，按照目录以及名字排序

	response.ProjectCatalog = catalog
	return response, nil
}
