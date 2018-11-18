package lang_worker

import (
	"be/common"
	"be/common/log"
	"be/worker/process"
	"path"
)

type PlainTextLangWorker struct {
	pm              *process.ProcessMgr
	codePath        string
	projectFullName string
	rawConfig       string
}

func (w *PlainTextLangWorker) Init(pm *process.ProcessMgr, codePath string, projectFullName string, rawConfig string) {
	w.pm = pm
	w.codePath = codePath
	w.projectFullName = projectFullName
	w.rawConfig = rawConfig
}

func (w *PlainTextLangWorker) OpenFile(p *FilePath) (*File, error) {
	result := &File{}
	result.Name = p.Name

	targetFile := path.Join(w.codePath, p.Path, p.Name)
	log.Debugf("打开文件 %s", targetFile)

	rawContent, err := common.GetFileContent(targetFile)
	if err != nil {
		return nil, err
	}
	result.RawContent = rawContent

	return result, nil
}

func (w *PlainTextLangWorker) GoToDefinition(p *Position) []*Position {
	return nil
}

func (w *PlainTextLangWorker) GoToTypeDefition(p *Position) []*Position {
	return nil
}

func (w *PlainTextLangWorker) GoToImplementation(p *Position) []*Position {
	return nil
}

func (w *PlainTextLangWorker) Close() {

}
