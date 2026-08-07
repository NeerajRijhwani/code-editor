package editor

import (
	"github.com/gdamore/tcell/v3"
)

func (e *Editor) HandleMouse(ev *tcell.EventMouse) {
	y, x := ev.Position()
	buttons := ev.Buttons()
	switch buttons {
	case tcell.ButtonPrimary:
		e.MouseClick(x, y)
	case tcell.WheelUp:
		e.MouseWheelUp()
	case tcell.WheelDown:
		e.MouseWheelDown()
	}
}

func (e *Editor) MouseWheelDown() {
	count := e.Buffer.LineCount()
	e.Renderer.IncreaseFirstLine(count)
	curx, cury := e.Cursor.Position()
	if count-curx > 1 && count-curx < 30 {
		e.Cursor.MoveDown()
	} else if curx-e.Renderer.FirstLine < 10 {
		e.Cursor.SetCursor(e.Renderer.FirstLine+10, cury)
	}

}

func (e *Editor) MouseWheelUp() {
	e.Renderer.DecreaseFirstLine()
	curx, cury := e.Cursor.Position()
	if curx-e.Renderer.FirstLine > 20 {
		e.Cursor.SetCursor(e.Renderer.FirstLine+20, cury)
	}

}

func (e *Editor) MouseClick(x, y int) {
	x, y = x-e.Renderer.Xoffset+e.Renderer.FirstLine, y-e.Renderer.Yoffset
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
