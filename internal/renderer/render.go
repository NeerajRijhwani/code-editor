package renderer

import (
	"log"
	"strconv"

	"github.com/gdamore/tcell/v3"
)

type Renderer struct {
	Screen    tcell.Screen
	Xoffset   int
	Yoffset   int
	FirstLine int
}

func (r *Renderer) SetClipboard(text string) {
	r.Screen.SetClipboard([]byte(text))
}

func (r *Renderer) Clear() {
	r.Screen.Clear()
}

func (r *Renderer) offset(x, y int) (int, int) {
	return x + r.Xoffset, y + r.Yoffset
}

func (r *Renderer) DecreaseFirstLine() {
	r.FirstLine = max(0, r.FirstLine-1)
}
func (r *Renderer) IncreaseFirstLine(count int) {
	log.Printf("First Line: %d and linecount : %d", r.FirstLine, count)
	if count-r.FirstLine > 30 {
		r.FirstLine = r.FirstLine + 1
	}
}

func (r *Renderer) DrawBox(x1, y1, x2, y2 int, style tcell.Style) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			r.Screen.Put(col, row, " ", style)
		}
	}

}

func (r *Renderer) SetActivelinestyle(x, width int, style tcell.Style) {
	offsetx, offsety := r.offset(0, 0)
	for i := offsety; i < width; i++ {
		r.Screen.SetContent(i, x+offsetx-r.FirstLine, ' ', nil, style)
	}
}

func (r *Renderer) DrawCell(x int, y int, ch rune, style tcell.Style) {
	X, Y := r.offset(x, y)
	r.Screen.SetContent(Y, X, ch, nil, style)
}

func (r *Renderer) DrawBorder(width, height int, style tcell.Style) {
	for col := 0; col <= width; col++ {
		r.Screen.Put(col, 0, string(tcell.RuneHLine), style)
		r.Screen.Put(col, height, string(tcell.RuneHLine), style)
	}

	for row := 0; row < height; row++ {
		r.Screen.Put(0, row, string(tcell.RuneVLine), style)
		r.Screen.Put(width, row, string(tcell.RuneVLine), style)
	}

}

func (r *Renderer) DrawCursor(x, y int) {
	x -= r.FirstLine
	x, y = r.offset(x, y)
	r.Screen.ShowCursor(y, x)
}

func (r *Renderer) DrawlineNumber(i int, style tcell.Style) {
	num := strconv.Itoa(i + r.FirstLine)
	for j, ch := range num {
		r.Screen.SetContent(r.Yoffset+j-len(num)-1, i+r.Xoffset, ch, nil, style)
	}
}

func (r *Renderer) Show() {
	r.Screen.Show()
}

func (r *Renderer) Sync() {
	r.Screen.Sync()
}

func (r *Renderer) Quit() {
	// You have to catch panics in a defer, clean up, and
	// re-raise them - otherwise your application can
	// die without leaving any diagnostic trace.
	maybePanic := recover()
	r.Screen.Fini()
	if maybePanic != nil {
		panic(maybePanic)
	}
}

func (r *Renderer) CursorStyleSet(mode rune) {
	switch mode {
	case 'v':
		fallthrough
	case 'n':
		r.Screen.SetCursorStyle(tcell.CursorStyleSteadyBlock)
	case 'i':
		r.Screen.SetCursorStyle(tcell.CursorStyleSteadyBar)
	default:
		r.Screen.SetCursorStyle(tcell.CursorStyleBlinkingBar)
	}
}

func InitRenderer() (*Renderer, error) {
	defStyle := tcell.StyleDefault.Background(tcell.GetColor("#1E1E2E"))

	s, err := tcell.NewScreen()

	if err != nil {
		return nil, err
	}

	if err := s.Init(); err != nil {
		return nil, err
	}

	s.SetStyle(defStyle)
	s.SetCursorStyle(tcell.CursorStyleSteadyBlock)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()

	return &Renderer{
		Screen:    s,
		Xoffset:   1,
		Yoffset:   5,
		FirstLine: 0,
	}, nil
}
