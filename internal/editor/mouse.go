package editor

import (
	"log"

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
	if curx-e.Renderer.FirstLine > 19 {
		e.Cursor.SetCursor(e.Renderer.FirstLine+20, cury)
	}

}

func (e *Editor) MouseClick(x, y int) {

	newx, newy := x-e.Renderer.Xoffset+e.Renderer.FirstLine, y-e.Renderer.Yoffset
	totalcols, err := e.Buffer.LineLength(x)
	if err != nil {
		return
	}
	log.Printf("First Line %d and x is %d and ", e.Renderer.FirstLine, x)

	if x-e.Renderer.Xoffset < 10 {
		e.Renderer.FirstLine = max(0, e.Renderer.FirstLine-(10-x+e.Renderer.Xoffset))
	}
	if x-e.Renderer.Xoffset > 19 {
		log.Printf("ans %d", x-e.Renderer.Xoffset)
		e.Renderer.FirstLine += x - e.Renderer.Xoffset - 19
		e.Renderer.FirstLine = min(e.Buffer.LineCount()-30, e.Renderer.FirstLine)
	}
	if newy > totalcols {
		e.Cursor.SetCursor(newx, totalcols-1)
	} else {
		e.Cursor.SetCursor(newx, newy)
	}

}
