package editor

import (
	"github.com/NeerajRijhwani/code-editor/internal/buffer"
	"github.com/NeerajRijhwani/code-editor/internal/cursor"
	"github.com/NeerajRijhwani/code-editor/internal/renderer"
	"github.com/gdamore/tcell/v3"
)

type Editor struct {
	Buffer   *buffer.Buffer
	Cursor   *cursor.Cursor
	Renderer *renderer.Renderer
	running  bool
}

func quitEditor(e *Editor) {
	e.Renderer.Quit()
}

func InitEditor() (*Editor, error) {
	b := buffer.InitBuffer()
	c := cursor.InitCursor()
	r, err := renderer.InitRenderer()
	if err != nil {
		return nil, err
	}
	return &Editor{
		Buffer:   b,
		Cursor:   c,
		Renderer: r,
		running:  true,
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

func (e *Editor) render() {
	e.Renderer.Clear()

	e.Renderer.DrawBorder(150, 150)

	e.Renderer.DrawBuffer(e.Buffer)

	e.Renderer.DrawCursor(e.Cursor)

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
			e.render()
			// case *tcell.EventMouse:
			// 	e.HandleMouse(ev)
		}
	}
	// fmt.Println("terminal has ended")
	defer quitEditor(e)
}
