package common

import "strings"

type LangType int64

const (
	// 无格式文件
	PlainText LangType = iota
	// Python
	Python
	// Go
	Go
)

// 映射关系，key为LangType，value为对应的文件后缀
var LangTypeMapping = map[LangType][]string{
	Python: []string{"py"},
	Go: []string{"go"},
}

func GetLangTypeByFileName(fileName string) LangType {
	fileName = strings.ToUpper(strings.TrimSpace(fileName))
	if fileName == "" {
		return PlainText
	}
	suffixes := strings.Split(fileName, ".")
	suffix := suffixes[len(suffixes)-1]

	for lt, s := range LangTypeMapping {
		for _, l := range s {
			if strings.ToUpper(strings.TrimSpace(suffix)) == strings.ToUpper(strings.TrimSpace(l)) {
				return lt
			}
		}
	}
	return PlainText
}

func GetLangTypeName(lt LangType) string {
	if suffixes, ok := LangTypeMapping[lt]; ok == true {
		return strings.Join(suffixes, ", ")
	} else {
		return strings.Join(LangTypeMapping[PlainText], ", ")
	}
}
