package lang_worker

import (
	"be/common"
	"be/common/log"
	"be/lex"
	"be/worker/process"
	"path"
)

type PythonLangWorker struct {
	pm              *process.ProcessMgr
	codePath        string
	projectFullName string
	rawConfig       string
}

func (w *PythonLangWorker) Init(pm *process.ProcessMgr, codePath string, projectFullName string, rawConfig string) {
	w.pm = pm
	w.codePath = codePath
	w.projectFullName = projectFullName
	w.rawConfig = rawConfig
}

func (w *PythonLangWorker) OpenFile(p *FilePath) (*File, error) {
	result := &File{}
	result.Name = p.Name

	targetFile := path.Join(w.codePath, p.Path, p.Name)

	tokens, err := lex.PygmentsLex(targetFile)
	if err != nil {
		log.Errorf("lex失败 %s: %s", targetFile, err.Error())
		// 如果失败则获取文件原始内容
		rawContent, err := common.GetFileContent(targetFile)
		if err != nil {
			return nil, err
		}
		result.RawContent = rawContent
	} else {
		result.Tokens = tokens
	}

	return result, nil
}

func (w *PythonLangWorker) GoToDefinition(p *Position) []*Position {
	return nil
}

func (w *PythonLangWorker) GoToTypeDefition(p *Position) []*Position {
	return nil
}

func (w *PythonLangWorker) GoToImplementation(p *Position) []*Position {
	return nil
}

func (w *PythonLangWorker) Close() {

}
