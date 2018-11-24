package lang_worker

import (
	"be/common/log"
	"be/lex"
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
	result := &File{Tokens: nil, RawContent: ""}
	result.Name = p.Name

	targetFile := path.Join(w.codePath, p.Path, p.Name)

	tokens, err := lex.PygmentsLex(targetFile)
	if err != nil {
		log.Errorf("lex失败 %s: %s", targetFile, err.Error())
		return nil, err
	} else {
		result.Tokens = tokens
	}

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
