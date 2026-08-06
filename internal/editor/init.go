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
)

type Editor struct {
	Buffer   *buffer.Buffer
	Cursor   *cursor.Cursor
	Renderer *renderer.Renderer
	Select   *cursor.Selection
	History  *buffer.Manager
	Status   *plugins.Status
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
	return &Editor{
		Buffer:   b,
		Cursor:   c,
		Renderer: r,
		Select:   s,
		History:  h,
		Status:   status,
		Mode:     'n',
		running:  true,
		FilePath: path,
		Modified: false,
		Height:   240,
		Width:    240,
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
	case tcell.WheelDown:
		log.Println("Wheel Down")
		e.Renderer.IncreaseFirstLine(count)
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
		line, _ := e.Buffer.GetLine(i + e.Renderer.FirstLine)

		for j, ch := range line {
			if e.Select.Active && e.Select.CheckWithinSelect(i, j) {
				e.Renderer.DrawCell(i, j, ch, 's')
			} else if i == cx-e.Renderer.FirstLine {
				e.Renderer.DrawCell(i, j, ch, 'a')
			} else {
				e.Renderer.DrawCell(i, j, ch, 'n')
			}
		}

		e.Renderer.DrawlineNumber(i)
	}

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
