package editor

import (
	"fmt"
	"log"

	"github.com/NeerajRijhwani/code-editor/internal/plugins"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func (e *Editor) drawCursor() {
	x, y := e.Cursor.Position()
	e.Renderer.DrawCursor(x, y)
}

func (e *Editor) reset() {
	e.Renderer.Clear()

	// Render Background
	e.Renderer.DrawBox(0, 0, e.Height, e.Width, e.Theme.Background)

	x, _ := e.Cursor.Position()
	// Render ActiveLine
	e.Renderer.SetActivelinestyle(x, e.Width, e.Theme.ActiveStyle)

	e.Renderer.CursorStyleSet(e.Mode)
}

func (e *Editor) Update() {
	e.Renderer.Show()
}

func (e *Editor) Render() {
	e.reset()

	e.drawBuffer()

	e.drawCursor()

	e.drawPlugins()

	e.Update()
}
func (e *Editor) drawPlugins() {
	e.RenderStatusLine()
}

func (t *Theme) GetColor(tokenText string) color.Color {
	switch tokenText {
	case "keyword":
		return t.KeywordStyle
	case "function":
		return t.FunctionStyle
	case "string":
		return t.StringStyle
	case "comment":
		return t.CommentStyle
	case "type":
		return t.TypeStyle
	case "number":
		return t.NumberStyle
	case "boolean", "constant":
		return t.ConstantStyle

	default:
		return color.White.TrueColor()
	}
}
func (e *Editor) drawBuffer() {

	count := e.Buffer.LineCount()
	cx, _ := e.Cursor.Position()
	tokens, err := e.TreeSitter.HighlightBuffer()
	if err != nil {
		log.Printf("Unable to highlight")
	}
	for i := range min(30, count-e.Renderer.FirstLine) {
		bufline := i + e.Renderer.FirstLine
		line, err := e.Buffer.GetLine(bufline)
		if err != nil {
			log.Printf("%v", err)
			return
		}
		cellcolor := color.White.TrueColor()
		lineTokens := e.getTokensForLine(tokens, uint32(bufline))
		for j, ch := range line {
			if lineTokens != nil {
				col := uint32(j)

				var bestToken *plugins.HighlightToken
				smallestRange := uint32(1<<32 - 1)

				for idx := range lineTokens {
					token := &lineTokens[idx]
					if col >= token.StartCol && col < token.EndCol {
						tokenRange := token.EndCol - token.StartCol

						if tokenRange < smallestRange {
							smallestRange = tokenRange
							bestToken = token
						}
					}
				}

				if bestToken != nil {
					cellcolor = e.Theme.GetColor(bestToken.Capture)
				}
			}

			if e.Select.Active && e.Select.CheckWithinSelect(i+e.Renderer.FirstLine, j) {
				cellstyle := e.Theme.SelectionStyle.Foreground(cellcolor)
				e.Renderer.DrawCell(i, j, ch, cellstyle)
			} else if i == cx-e.Renderer.FirstLine {
				cellstyle := e.Theme.ActiveStyle.Foreground(cellcolor)
				e.Renderer.DrawCell(i, j, ch, cellstyle)
			} else {
				cellstyle := e.Theme.TextStyle.Foreground(cellcolor)
				e.Renderer.DrawCell(i, j, ch, cellstyle)
			}
		}

		e.Renderer.DrawlineNumber(i, e.Theme.LinenumStyle)
	}
}

func (e *Editor) getTokensForLine(tokens []plugins.HighlightToken, lineNum uint32) []plugins.HighlightToken {
	if tokens == nil {
		return nil
	}
	var lineTokens []plugins.HighlightToken
	for _, token := range tokens {
		if token.StartLine <= lineNum && token.EndLine >= lineNum {
			lineTokens = append(lineTokens, token)
		}
	}
	return lineTokens
}

func (e *Editor) RenderStatusLine() {
	x, y := e.Cursor.Position()
	mode := ""
	var modestyle tcell.Style
	switch e.Mode {
	case 'n':
		mode = " NORMAL "
		modestyle = e.Theme.NormalModeStyle
	case 'i':
		mode = " INSERT "
		modestyle = e.Theme.InsertModeStyle
	case 'v':
		mode = " VISUAL "
		modestyle = e.Theme.VisualModeStyle
	}

	background := tcell.StyleDefault.Background(tcell.GetColor("#303030"))
	startPos := 31
	height := 2
	width := 152
	style := background.Foreground(color.Gray.TrueColor())

	e.Renderer.DrawBox(0, startPos+height/2, width, startPos+height/2, background)
	pos := 0
	for _, ch := range mode {
		e.Renderer.Screen.SetContent(pos, startPos+(height/2), ch, nil, modestyle)
		pos++
	}
	filename := " " + e.FileName + " "
	for _, ch := range filename {
		e.Renderer.Screen.SetContent(pos, startPos+(height/2), ch, nil, style)
		pos++
	}

	if !e.Modified {
		for _, ch := range "[+]" {
			e.Renderer.Screen.SetContent(pos, startPos+(height/2), ch, nil, style)
			pos++
		}
	}
	cursortext := fmt.Sprintf(" %d:%d ", x, y)
	pos = width
	for i := len(cursortext) - 1; i >= 0; i-- {
		e.Renderer.Screen.SetContent(pos, startPos+(height/2), rune(cursortext[i]), nil, modestyle)
		pos--
	}
	filetype := " " + e.FileType + " "
	for i := len(filetype) - 1; i >= 0; i-- {
		e.Renderer.Screen.SetContent(pos, startPos+(height/2), rune(filetype[i]), nil, style)
		pos--
	}

}
