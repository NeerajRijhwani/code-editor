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

type TreeSitter struct {
	lang   *sitter.Language
	tree   *sitter.Tree
	parser *sitter.Parser
	query  *sitter.Query
}

func loadquery(lang string) ([]byte, error) {
	if lang == "go" {
		content, err := os.ReadFile("./internal/plugins/highlights_go.scm")
		if err != nil {
			return nil, err
		}
		return content, nil
	}
	return nil, nil
}

func (t *TreeSitter) CleanUp() {
	t.tree.Close()
	t.query.Close()
	t.parser.Close()
}

func Init_Treesitter(sourcecode []byte) (*TreeSitter, error) {
	parser := sitter.NewParser()
	lang := golang.GetLanguage()
	parser.SetLanguage(lang)

	tree, err := parser.ParseCtx(context.Background(), nil, sourcecode)
	if err != nil {
		log.Println("Parser failed")
		return nil, err
	}
	highlightquery, err := loadquery("go")

	if err != nil {
		log.Println("Unable to Query")
		return nil, err
	}
	query, err := sitter.NewQuery(highlightquery, lang)
	if err != nil {
		return nil, fmt.Errorf("failed to compile query: %w", err)
	}
	return &TreeSitter{
		lang:   lang,
		parser: parser,
		tree:   tree,
		query:  query,
	}, nil
}

func (t *TreeSitter) HighlightBuffer() ([]HighlightToken, error) {
	qc := sitter.NewQueryCursor()
	defer qc.Close()

	rootNode := t.tree.RootNode()
	qc.Exec(t.query, rootNode)

	var tokens []HighlightToken
	for {
		match, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, capture := range match.Captures {
			captureName := t.query.CaptureNameForId(capture.Index)
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

func GetEdit(startidx, oldendidx, newendidx, StartPosx, StartPosy, oldendposx, oldendposy, newendposx, newendposy uint32) sitter.EditInput {
	return sitter.EditInput{
		StartIndex:  startidx,
		OldEndIndex: oldendidx,
		NewEndIndex: newendidx,
		StartPoint: sitter.Point{
			Row:    StartPosx,
			Column: StartPosy,
		},
		OldEndPoint: sitter.Point{
			Row:    oldendposx,
			Column: oldendposy,
		},
		NewEndPoint: sitter.Point{
			Row:    newendposx,
			Column: newendposy,
		},
	}
}

func (t *TreeSitter) ApplyEdit(source []byte, edit sitter.EditInput) error {
	t.tree.Edit(edit)
	newTree, err := t.parser.ParseCtx(
		context.Background(),
		t.tree,
		source,
	)
	if err != nil {
		return err
	}

	t.tree.Close()
	t.tree = newTree

	return nil
}
