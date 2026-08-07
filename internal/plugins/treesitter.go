package plugins

import (
	"context"
	"fmt"
	"log"
	"os"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
)

type HighlightToken struct {
	StartLine uint32
	StartCol  uint32
	EndLine   uint32
	EndCol    uint32
	Capture   string
}

func loadquery(lang string) ([]byte, error) {
	if lang == "go" {
		content, err := os.ReadFile("./internal/plugins/highlishts_go.scm")
		if err != nil {
			return nil, err
		}
		return content, nil
	}
	return nil, nil
}

func HighlightBuffer(sourcecode []byte) ([]HighlightToken, error) {
	parser := sitter.NewParser()
	lang := golang.GetLanguage()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(context.Background(), nil, sourcecode)
	if err != nil {
		log.Println("Parser failed")
		return nil, err
	}
	defer tree.Close()

	rootNode := tree.RootNode()
	highlightquery, err := loadquery("go")

	if err != nil {
		log.Println("Unable to Query")
		return nil, err
	}
	query, err := sitter.NewQuery(highlightquery, lang)
	if err != nil {
		return nil, fmt.Errorf("failed to compile query: %w", err)
	}
	defer query.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	qc.Exec(query, rootNode)

	var tokens []HighlightToken
	for {
		match, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			captureName := query.CaptureNameForId(capture.Index)
			node := capture.Node

			tokens = append(tokens, HighlightToken{
				StartLine: node.StartPoint().Row,
				StartCol:  node.StartPoint().Column,
				EndLine:   node.EndPoint().Row,
				EndCol:    node.EndPoint().Column,
				Capture:   captureName,
			})
		}
	}

	return tokens, nil
}
