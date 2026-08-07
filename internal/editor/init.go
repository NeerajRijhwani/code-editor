package editor

import (
	"log"
	"path/filepath"
	"unicode"

	"github.com/NeerajRijhwani/code-editor/internal/buffer"
	"github.com/NeerajRijhwani/code-editor/internal/cursor"
	"github.com/NeerajRijhwani/code-editor/internal/plugins"
	"github.com/NeerajRijhwani/code-editor/internal/renderer"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

func (t *Theme) GetHashedColor(tokenText string) color.Color {
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

type Theme struct {
	Background     tcell.Style
	TextStyle      tcell.Style
	BorderStyle    tcell.Style
	LinenumStyle   tcell.Style
	ActiveStyle    tcell.Style
	CursorStyle    tcell.CursorStyle
	StatusStyle    tcell.Style
	SelectionStyle tcell.Style
	KeywordStyle   color.Color
	FunctionStyle  color.Color
	StringStyle    color.Color
	CommentStyle   color.Color
	TypeStyle      color.Color
	NumberStyle    color.Color
	ConstantStyle  color.Color
}

type Editor struct {
	Buffer    *buffer.Buffer
	Cursor    *cursor.Cursor
	Renderer  *renderer.Renderer
	Select    *cursor.Selection
	History   *buffer.Manager
	Status    *plugins.Status
	Mode      rune
	Width     int
	Height    int
	FilePath  string
	Modified  bool
	running   bool
	Theme     Theme
	Highlight []plugins.HighlightToken
}

func quitEditor(e *Editor) {
	e.Renderer.Quit()
}

func InitEditor(path string) (*Editor, error) {
	b, err := OpenFile(path)
	if err != nil {
		log.Println("Cannot inialize buffer")
		return nil, err
	}
	c := cursor.InitCursor()
	r, err := renderer.InitRenderer()
	s := cursor.InitSelect()
	h := buffer.Init_Manager()
	filetype := filepath.Ext(path)
	runes := []rune(filetype[1:])
	runes[0] = unicode.ToUpper(runes[0])

	status := plugins.Init_StatusLine(2, 152, 31, filepath.Base(path), string(runes))
	if err != nil {
		return nil, err
	}
	Highlight, err := plugins.HighlightBuffer(b.GetBuffer())
	if err != nil {
		return nil, err
	}
	return &Editor{
		Buffer:    b,
		Cursor:    c,
		Renderer:  r,
		Select:    s,
		History:   h,
		Status:    status,
		Mode:      'n',
		running:   true,
		FilePath:  path,
		Modified:  false,
		Height:    240,
		Width:     240,
		Highlight: Highlight,
		Theme: Theme{
			Background: tcell.StyleDefault.
				Background(tcell.GetColor("#050505")),

			TextStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#050505")),

			BorderStyle: tcell.StyleDefault.
				Foreground(tcell.GetColor("#6C7086")),

			LinenumStyle: tcell.StyleDefault.
				Foreground(tcell.GetColor("#7F849C")).Background(tcell.GetColor("#050505")),

			ActiveStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#37373b")),

			CursorStyle: tcell.CursorStyleSteadyBlock,

			StatusStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#45475A")).Foreground(tcell.GetColor("#CDD6F4")),

			SelectionStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#9c9a97")),

			KeywordStyle:  color.NewRGBColor(203, 166, 247),
			FunctionStyle: color.NewRGBColor(137, 180, 250),
			StringStyle:   color.NewRGBColor(166, 227, 161),
			CommentStyle:  color.NewRGBColor(108, 112, 134),
			TypeStyle:     color.NewRGBColor(249, 226, 175),
			NumberStyle:   color.NewRGBColor(250, 179, 135),
			ConstantStyle: color.NewRGBColor(243, 139, 168),
		},
	}, nil
}

func (e *Editor) HandleMouse(ev *tcell.EventMouse) {
	y, x := ev.Position()
	count := e.Buffer.LineCount()
	buttons := ev.Buttons()
	switch buttons {
	case tcell.ButtonPrimary:
		e.MouseClick(x, y)
	case tcell.WheelUp:
		log.Println("Wheel Up")
		e.Renderer.DecreaseFirstLine()
		curx, cury := e.Cursor.Position()
		if curx-e.Renderer.FirstLine > 20 {
			e.Cursor.SetCursor(e.Renderer.FirstLine+20, cury)
		}

	case tcell.WheelDown:
		log.Println("Wheel Down")
		e.Renderer.IncreaseFirstLine(count)
		curx, cury := e.Cursor.Position()
		if count-curx > 1 && count-curx < 30 {
			e.Cursor.MoveDown()
		} else if curx-e.Renderer.FirstLine < 10 {
			e.Cursor.SetCursor(e.Renderer.FirstLine+10, cury)
		}
	}
}

func (e *Editor) MouseClick(x, y int) {
	x, y = x-e.Renderer.Xoffset+e.Renderer.FirstLine, y-e.Renderer.Yoffset
	log.Printf("Mouse Clicked At %d and %d ", x, y)
	totalcols, err := e.Buffer.LineLength(x)
	if err != nil {
		return
	}
	if y > totalcols {
		e.Cursor.SetCursor(x, totalcols-1)
	} else {
		e.Cursor.SetCursor(x, y)
	}

}

func (e *Editor) drawBuffer() {

	count := e.Buffer.LineCount()
	cx, _ := e.Cursor.Position()
	for i := range min(30, count-e.Renderer.FirstLine) {
		bufline := i + e.Renderer.FirstLine
		line, _ := e.Buffer.GetLine(bufline)
		lineTokens := e.getTokensForLine(uint32(bufline))
		for j, ch := range line {
			col := uint32(j)
			cellcolor := color.White.TrueColor()

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
				cellcolor = e.Theme.GetHashedColor(bestToken.Capture)
			}

			if e.Select.Active && e.Select.CheckWithinSelect(i, j) {
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

		e.Renderer.DrawlineNumber(i)
	}

}

func (e *Editor) getTokensForLine(lineNum uint32) []plugins.HighlightToken {
	var lineTokens []plugins.HighlightToken
	for _, token := range e.Highlight {
		if token.StartLine <= lineNum && token.EndLine >= lineNum {
			lineTokens = append(lineTokens, token)
		}
	}
	return lineTokens
}

func (e *Editor) drawCursor() {
	x, y := e.Cursor.Position()
	e.Renderer.DrawCursor(x, y)
}

func (e *Editor) reset() {
	e.Renderer.Clear()

	// Render Background
	e.Renderer.DrawBox(0, 0, e.Height, e.Width, e.Renderer.Theme.Background)

	x, _ := e.Cursor.Position()
	// Render ActiveLine
	e.Renderer.SetActivelinestyle(x, e.Width)

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
	x, y := e.Cursor.Position()
	e.Status.RenderStatusLine(e.Renderer, x, y)
}

func (e *Editor) Run() {
	// fmt.Println("terminal has started")
	for e.running {
		ev := <-e.Renderer.Screen.EventQ()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			e.Renderer.Sync()
		case *tcell.EventKey:
			e.HandleKey(ev)
			e.Render()
		case *tcell.EventMouse:
			e.HandleMouse(ev)
			e.Render()
		}
	}
	// fmt.Println("terminal has ended")
	defer quitEditor(e)
}
