package editor

import (
	"log"

	"github.com/NeerajRijhwani/code-editor/internal/buffer"
	"github.com/NeerajRijhwani/code-editor/internal/cursor"
	"github.com/NeerajRijhwani/code-editor/internal/renderer"
	"github.com/gdamore/tcell/v3"
)

type Editor struct {
	Buffer   *buffer.Buffer
	Cursor   *cursor.Cursor
	Renderer *renderer.Renderer
	Select   *cursor.Selection
	Mode     rune
	Width    int
	Height   int
	FilePath string
	Modified bool
	running  bool
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
	r, err := renderer.InitRenderer(b)
	s := cursor.InitSelect()
	if err != nil {
		return nil, err
	}
	return &Editor{
		Buffer:   b,
		Cursor:   c,
		Renderer: r,
		Select:   s,
		Mode:     'n',
		running:  true,
		FilePath: path,
		Modified: false,
		Height:   240,
		Width:    240,
	}, nil
}

func (e *Editor) HandleMouse(ev *tcell.EventMouse) {
	x, y := ev.Position()
	totalcols, _ := e.Buffer.LineLength(x)
	if y > totalcols {
		e.Cursor.SetCursor(x, totalcols-1)
	} else {
		e.Cursor.SetCursor(x, y)
	}
}

func (e *Editor) Render() {
	e.Renderer.Clear()

	// Render Background
	e.Renderer.DrawBox(0, 0, e.Height, e.Width, e.Renderer.Theme.Background)
	// e.Renderer.DrawBorder(170, 150)

	e.Renderer.DrawBuffer(e.Buffer, e.Cursor, e.Select)

	e.Renderer.CursorStyleSet(e.Mode)

	e.Renderer.DrawCursor(e.Cursor)

	// Render Selection

	e.Renderer.Show()
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
			// case *tcell.EventMouse:
			// 	e.HandleMouse(ev)
		}
	}
	// fmt.Println("terminal has ended")
	defer quitEditor(e)
}
