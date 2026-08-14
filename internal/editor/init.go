package editor

import (
	"log"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/NeerajRijhwani/code-editor/internal/buffer"
	"github.com/NeerajRijhwani/code-editor/internal/cursor"
	"github.com/NeerajRijhwani/code-editor/internal/plugins"
	"github.com/NeerajRijhwani/code-editor/internal/renderer"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

const (
	NORMAL = 'n'
	INSERT = 'i'
	VISUAL = 'v'
)

type Theme struct {
	Background      tcell.Style
	TextStyle       tcell.Style
	BorderStyle     tcell.Style
	LinenumStyle    tcell.Style
	ActiveStyle     tcell.Style
	NormalModeStyle tcell.Style
	InsertModeStyle tcell.Style
	VisualModeStyle tcell.Style
	StatusStyle     tcell.Style
	SelectionStyle  tcell.Style
	CursorStyle     tcell.CursorStyle
	KeywordStyle    color.Color
	FunctionStyle   color.Color
	StringStyle     color.Color
	CommentStyle    color.Color
	TypeStyle       color.Color
	NumberStyle     color.Color
	ConstantStyle   color.Color
}

type Editor struct {
	FileName            string
	FileType            string
	Pasting_Buffer      strings.Builder
	Pasting_In_Progress bool
	Buffer              *buffer.Buffer
	Cursor              *cursor.Cursor
	Renderer            *renderer.Renderer
	Select              *cursor.Selection
	History             *buffer.Manager
	Mode                rune
	Width               int
	Height              int
	FilePath            string
	Modified            bool
	running             bool
	Theme               Theme
	Highlight           []plugins.HighlightToken
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

	if err != nil {
		return nil, err
	}
	Highlight, err := plugins.HighlightBuffer(b.GetBuffer())
	if err != nil {
		return nil, err
	}
	var pastingbuf strings.Builder
	return &Editor{
		FileName:            filepath.Base(path),
		FileType:            filetype,
		FilePath:            path,
		Pasting_Buffer:      pastingbuf,
		Pasting_In_Progress: false,
		Buffer:              b,
		Cursor:              c,
		Renderer:            r,
		Select:              s,
		History:             h,
		Mode:                'n',
		running:             true,
		Modified:            false,
		Height:              240,
		Width:               240,
		Highlight:           Highlight,
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

			StatusStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#45475A")).Foreground(tcell.GetColor("#CDD6F4")),

			SelectionStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#9c9a97")),

			NormalModeStyle: tcell.StyleDefault.
				Background(color.IndianRed).Foreground(color.Black.TrueColor()),

			InsertModeStyle: tcell.StyleDefault.
				Background(color.DarkCyan).Foreground(color.Black.TrueColor()),

			VisualModeStyle: tcell.StyleDefault.
				Background(color.Orange).Foreground(color.Black.TrueColor()),

			CursorStyle: tcell.CursorStyleSteadyBlock,

			KeywordStyle: color.NewRGBColor(203, 166, 247),

			FunctionStyle: color.NewRGBColor(137, 180, 250),

			StringStyle: color.NewRGBColor(166, 227, 161),

			CommentStyle: color.NewRGBColor(108, 112, 134),

			TypeStyle: color.NewRGBColor(249, 226, 175),

			NumberStyle: color.NewRGBColor(250, 179, 135),

			ConstantStyle: color.NewRGBColor(243, 139, 168),
		},
	}, nil
}

func (e *Editor) HandlePaste(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyRune:
		e.Pasting_Buffer.WriteString(ev.Str())
	case tcell.KeyCtrlJ:
		e.Pasting_Buffer.WriteRune('\n')
	case tcell.KeyTab:
		e.Pasting_Buffer.WriteString("    ")
	}

}

func (e *Editor) Paste() {
	x, y := e.Cursor.Position()
	e.Buffer.InsertText(e.Pasting_Buffer.String(), x, y)
}

func (e *Editor) Run() {
	// fmt.Println("terminal has started")
	for e.running {
		ev := <-e.Renderer.Screen.EventQ()
		switch ev := ev.(type) {
		case *tcell.EventResize:
			e.Renderer.Sync()
		case *tcell.EventKey:
			if e.Pasting_In_Progress {
				e.HandlePaste(ev)
				continue
			}
			e.HandleKey(ev)
			e.Render()
		case *tcell.EventMouse:
			e.HandleMouse(ev)
			e.Render()
		case *tcell.EventPaste:
			if ev.Start() {
				log.Printf("pasting start")
				e.Pasting_In_Progress = true
			}
			if ev.End() {
				log.Printf("pasting end")
				e.Paste()
				e.Pasting_In_Progress = false
			}
		}
	}
	// fmt.Println("terminal has ended")
	defer quitEditor(e)
}
