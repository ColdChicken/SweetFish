package lang_worker

import "be/worker/process"

type PlainTextLangWorker struct {
}

func (w *PlainTextLangWorker) Init(pm *process.ProcessMgr, codePath string, projectFullName string, rawConfig string) {

}

func (w *PlainTextLangWorker) OpenFile(p *FilePath) *File {
	return nil
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
