package lex

import (
	"be/common"
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

func PygmentsLex(filename string) (*FileTokensInfo, error) {
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
				Type:    strings.Join(result.Type, "."),
				BegPos:  result.BegPos,
				EndPos:  result.EndPos,
			}
			resultLine = append(resultLine, token)
		}
		tokensInfo.Tokens = append(tokensInfo.Tokens, resultLine)
	}
	return tokensInfo, nil
}
