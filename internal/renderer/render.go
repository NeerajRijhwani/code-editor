package renderer

import (
	"log"
	"strconv"

	"github.com/NeerajRijhwani/code-editor/internal/buffer"
	"github.com/NeerajRijhwani/code-editor/internal/cursor"
	"github.com/gdamore/tcell/v3"
	// "github.com/gdamore/tcell/v3/color"
)

type Theme struct {
	Background     tcell.Style
	TextStyle      tcell.Style
	BorderStyle    tcell.Style
	LinenumStyle   tcell.Style
	ActiveStyle    tcell.Style
	CursorStyle    tcell.CursorStyle
	StatusStyle    tcell.Style
	SelectionStyle tcell.Style
}

type Renderer struct {
	Screen    tcell.Screen
	Theme     Theme
	Xoffset   int
	Yoffset   int
	LastLine  int
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
	r.FirstLine--
}
func (r *Renderer) IncreaseFirstLine() {
	r.FirstLine++
}
func (r *Renderer) DecreaseLastLine() {
	r.LastLine--
}
func (r *Renderer) IncreaseLastLine() {
	r.LastLine++
}

func (r *Renderer) UpdateLastLine(b *buffer.Buffer) {
	count := b.LineCount()
	if count < 30 {
		r.LastLine = count
	} else {
		r.LastLine = r.FirstLine + r.FirstLine%30
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
	// Only draw corners if necessary
	// if y1 != y2 && x1 != x2 {
	// 	r.Screen.Put(x1, y1, string(tcell.RuneULCorner), style)
	// 	r.Screen.Put(x2, y1, string(tcell.RuneURCorner), style)
	// 	r.Screen.Put(x1, y2, string(tcell.RuneLLCorner), style)
	// 	r.Screen.Put(x2, y2, string(tcell.RuneLRCorner), style)
	// }
}

func (r *Renderer) DrawText(x, y int, text string, c *cursor.Cursor, s *cursor.Selection) {
	x, y = r.offset(x, y)
	cx, _ := c.Position()
	for i, ch := range text {
		if s.Active && s.CheckWithinSelect(x-r.Xoffset, y+i-r.Yoffset) {
			log.Print("select is working")
			r.Screen.SetContent(y+i, x, ch, nil, r.Theme.SelectionStyle)
		} else if x == cx+r.Xoffset-r.FirstLine {
			r.Screen.SetContent(y+i, x, ch, nil, r.Theme.ActiveStyle)
		} else {
			r.Screen.SetContent(y+i, x, ch, nil, r.Theme.TextStyle)
		}
	}
}
func (r *Renderer) DrawBorder(width, height int) {
	for col := 0; col <= width; col++ {
		r.Screen.Put(col, 0, string(tcell.RuneHLine), r.Theme.BorderStyle)
		r.Screen.Put(col, height, string(tcell.RuneHLine), r.Theme.BorderStyle)
	}

	for row := 0; row < height; row++ {
		r.Screen.Put(0, row, string(tcell.RuneVLine), r.Theme.BorderStyle)
		r.Screen.Put(width, row, string(tcell.RuneVLine), r.Theme.BorderStyle)
	}

}

func (r *Renderer) DrawCursor(c *cursor.Cursor) {
	x, y := c.Position()
	x -= r.FirstLine
	x, y = r.offset(x, y)
	r.Screen.ShowCursor(y, x)
}

func (r *Renderer) DrawBuffer(b *buffer.Buffer, c *cursor.Cursor, s *cursor.Selection) {
	count := b.LineCount()
	for i := range min(30, count-r.FirstLine) {
		line, _ := b.GetLine(i + r.FirstLine)
		r.DrawText(i, 0, line, c, s)
		r.DrawlineNumber(i)
	}
}

func (r *Renderer) DrawlineNumber(i int) {
	num := strconv.Itoa(i + r.FirstLine)
	for j, ch := range num {
		r.Screen.SetContent(r.Yoffset+j-len(num)-1, i+r.Xoffset, ch, nil, r.Theme.LinenumStyle)
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

func InitRenderer(b *buffer.Buffer) (*Renderer, error) {
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
		Screen: s,
		Theme: Theme{
			Background: tcell.StyleDefault.
				Background(tcell.GetColor("#050505")),

			TextStyle: tcell.StyleDefault.
				Foreground(tcell.GetColor("#f0f2f0")).Background(tcell.GetColor("#050505")),

			BorderStyle: tcell.StyleDefault.
				Foreground(tcell.GetColor("#6C7086")),

			LinenumStyle: tcell.StyleDefault.
				Foreground(tcell.GetColor("#7F849C")).Background(tcell.GetColor("#050505")),

			ActiveStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#37373b")).Foreground(tcell.GetColor("#f0f2f0")),

			CursorStyle: tcell.CursorStyleSteadyBlock,
			StatusStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#45475A")).
				Foreground(tcell.GetColor("#CDD6F4")),
			SelectionStyle: tcell.StyleDefault.
				Background(tcell.GetColor("#9c9a97")).Foreground(tcell.GetColor("#f0f2f0")),
		},
		Xoffset:   1,
		Yoffset:   5,
		LastLine:  min(30, b.LineCount()),
		FirstLine: 0,
	}, nil
}

// func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
// 	if y2 < y1 {
// 		y1, y2 = y2, y1
// 	}
// 	if x2 < x1 {
// 		x1, x2 = x2, x1
// 	}
//
// 	// Fill background
// 	for row := y1; row <= y2; row++ {
// 		for col := x1; col <= x2; col++ {
// 			s.Put(col, row, " ", style)
// 		}
// 	}
//
// 	// Draw borders
// 	for col := x1; col <= x2; col++ {
// 		s.Put(col, y1, string(tcell.RuneHLine), style)
// 		s.Put(col, y2, string(tcell.RuneHLine), style)
// 	}
// 	for row := y1 + 1; row < y2; row++ {
// 		s.Put(x1, row, string(tcell.RuneVLine), style)
// 		s.Put(x2, row, string(tcell.RuneVLine), style)
// 	}
//
// 	// Only draw corners if necessary
// 	if y1 != y2 && x1 != x2 {
// 		s.Put(x1, y1, string(tcell.RuneULCorner), style)
// 		s.Put(x2, y1, string(tcell.RuneURCorner), style)
// 		s.Put(x1, y2, string(tcell.RuneLLCorner), style)
// 		s.Put(x2, y2, string(tcell.RuneLRCorner), style)
// 	}
//
// 	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
// }
// func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
// 	row := y1
// 	col := x1
// 	var width int
// 	for text != "" {
// 		text, width = s.Put(col, row, text, style)
// 		col += width
// 		if col >= x2 {
// 			row++
// 			col = x1
// 		}
// 		if row > y2 {
// 			break
// 		}
// 		if width == 0 {
// 			// incomplete grapheme at end of string
// 			break
// 		}
// 	}
// }
