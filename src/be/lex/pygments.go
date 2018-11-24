package lex

import (
	"be/common"
	"be/common/log"
	"be/options"
	"strings"
)

type Token struct {
	Content string `json:"content"`
	Type    string `json:"type"`
	BegPos  int64  `json:"beg_pos"`
	EndPos  int64  `json:"end_pos"`
}

type FileTokensInfo struct {
	Tokens [][]*Token `json:"tokens"`
}

// https://bitbucket.org/birkenfeld/pygments-main/src/7941677dc77d4f2bf0bbd6140ade85a9454b8b80/pygments/token.py?at=default&fileviewer=file-view-default
var tokenMap = map[string]string{
	"Token":                  "",
	"Text":                   "",
	"Whitespace":             "w",
	"Escape":                 "esc",
	"Error":                  "err",
	"Other":                  "x",
	"Keyword":                "k",
	"Keyword.Constant":       "kc",
	"Keyword.Declaration":    "kd",
	"Keyword.Namespace":      "kn",
	"Keyword.Pseudo":         "kp",
	"Keyword.Reserved":       "kr",
	"Keyword.Type":           "kt",
	"Name":                   "n",
	"Name.Attribute":         "na",
	"Name.Builtin":           "nb",
	"Name.Builtin.Pseudo":    "bp",
	"Name.Class":             "nc",
	"Name.Constant":          "no",
	"Name.Decorator":         "nd",
	"Name.Entity":            "ni",
	"Name.Exception":         "ne",
	"Name.Function":          "nf",
	"Name.Function.Magic":    "fm",
	"Name.Property":          "py",
	"Name.Label":             "nl",
	"Name.Namespace":         "nn",
	"Name.Other":             "nx",
	"Name.Tag":               "nt",
	"Name.Variable":          "nv",
	"Name.Variable.Class":    "vc",
	"Name.Variable.Global":   "vg",
	"Name.Variable.Instance": "vi",
	"Name.Variable.Magic":    "vm",
	"Literal":                "l",
	"Literal.Date":           "ld",
	"String":                 "s",
	"String.Affix":           "sa",
	"String.Backtick":        "sb",
	"String.Char":            "sc",
	"String.Delimiter":       "dl",
	"String.Doc":             "sd",
	"String.Double":          "s2",
	"String.Escape":          "se",
	"String.Heredoc":         "sh",
	"String.Interpol":        "si",
	"String.Other":           "sx",
	"String.Regex":           "sr",
	"String.Single":          "s1",
	"String.Symbol":          "ss",
	"Number":                 "m",
	"Number.Bin":             "mb",
	"Number.Float":           "mf",
	"Number.Hex":             "mh",
	"Number.Integer":         "mi",
	"Number.Integer.Long":    "il",
	"Number.Oct":             "mo",
	"Operator":               "o",
	"Operator.Word":          "ow",
	"Punctuation":            "p",
	"Comment":                "c",
	"Comment.Hashbang":       "ch",
	"Comment.Multiline":      "cm",
	"Comment.Preproc":        "cp",
	"Comment.PreprocFile":    "cpf",
	"Comment.Single":         "c1",
	"Comment.Special":        "cs",
	"Generic":                "g",
	"Generic.Deleted":        "gd",
	"Generic.Emph":           "ge",
	"Generic.Error":          "gr",
	"Generic.Heading":        "gh",
	"Generic.Inserted":       "gi",
	"Generic.Output":         "go",
	"Generic.Prompt":         "gp",
	"Generic.Strong":         "gs",
	"Generic.Subheading":     "gu",
	"Generic.Traceback":      "gt",
}

func getTokenType(types []string) string {
	if target, ok := tokenMap[strings.Join(types, ".")]; ok {
		return target
	} else {
		log.Warnf("token类型不存在缩写信息 %s", strings.Join(types, "."))
		if len(types) != 1 {
			return getTokenType(types[0 : len(types)-1])
		} else {
			return types[0]
		}
	}
}

func PygmentsLex(filename string) (*FileTokensInfo, error) {
	result, err := pygmentsLex(filename)
	if err == nil {
		return result, nil
	} else {
		rawContent, err := common.GetFileContent(filename)
		if err != nil {
			return nil, err
		} else {
			result = &FileTokensInfo{Tokens: [][]*Token{}}
			for _, contentLine := range strings.Split(rawContent, "\n") {
				line := []*Token{}
				token := &Token{
					Content: contentLine,
					Type:    "",
					BegPos:  0,
					EndPos:  int64(len(contentLine)),
				}
				line = append(line, token)
				result.Tokens = append(result.Tokens, line)
			}
			return result, nil
		}
	}
}

func pygmentsLex(filename string) (*FileTokensInfo, error) {
	args := []string{options.Options.PygmentsHelperPath, filename}
	o, _, err := common.Exec(30, "python", args...)
	if err != nil {
		return nil, err
	}

	type pygmentsOutput struct {
		Content string   `json:"content"`
		Type    []string `json:"type"`
		BegPos  int64    `json:"beg_pos"`
		EndPos  int64    `json:"end_pos"`
	}

	results := [][]*pygmentsOutput{}
	if err := common.ParseJsonStr(o, &results); err != nil {
		return nil, err
	}

	tokensInfo := &FileTokensInfo{Tokens: [][]*Token{}}

	for _, line := range results {
		resultLine := []*Token{}
		for _, result := range line {
			token := &Token{
				Content: result.Content,
				Type:    getTokenType(result.Type),
				BegPos:  result.BegPos,
				EndPos:  result.EndPos,
			}
			resultLine = append(resultLine, token)
		}
		tokensInfo.Tokens = append(tokensInfo.Tokens, resultLine)
	}
	return tokensInfo, nil
}
