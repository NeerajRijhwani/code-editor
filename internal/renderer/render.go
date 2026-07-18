package renderer

import (
	"github.com/NeerajRijhwani/code-editor/internal/buffer"
	"github.com/NeerajRijhwani/code-editor/internal/cursor"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

type Renderer struct {
	Screen      tcell.Screen
	screenstyle tcell.Style
	cellstyle   tcell.Style
	BorderStyle tcell.Style
	Xoffset     int
	Yoffset     int
}

func (r *Renderer) Clear() {
	r.Screen.Clear()
}

func (r *Renderer) offset(x, y int) (int, int) {
	return x + r.Xoffset, y + r.Yoffset
}
func (r *Renderer) DrawText(x, y int, text string) {
	x, y = r.offset(x, y)
	for i, ch := range text {
		r.Screen.SetContent(x+i, y, ch, nil, r.cellstyle)
	}
}
func (r *Renderer) DrawBorder(width, height int) {
	for col := 0; col <= width; col++ {
		r.Screen.Put(col, 0, string(tcell.RuneHLine), r.BorderStyle)
		r.Screen.Put(col, height, string(tcell.RuneHLine), r.BorderStyle)
	}

	for row := 0; row < height; row++ {
		r.Screen.Put(0, row, string(tcell.RuneVLine), r.BorderStyle)
		r.Screen.Put(width, row, string(tcell.RuneVLine), r.BorderStyle)
	}

}

func (r *Renderer) DrawCursor(c *cursor.Cursor) {
	x, y := c.Position()
	x, y = r.offset(x, y)
	r.Screen.ShowCursor(y, x)
}

func (r *Renderer) DrawBuffer(b *buffer.Buffer) {
	count := b.LineCount()
	for i := range count {
		line, _ := b.GetLine(i)
		r.DrawText(0, i, line)
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

func InitRenderer() (*Renderer, error) {
	defStyle := tcell.StyleDefault.Background(color.Black).Foreground(color.White)
	s, err := tcell.NewScreen()

	if err != nil {
		return nil, err
	}

	if err := s.Init(); err != nil {
		return nil, err
	}

	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	s.Clear()
	cellstyle := tcell.StyleDefault.Foreground(color.White).Background(color.Black)
	borderstyle := tcell.StyleDefault.Foreground(color.White).Background(color.Black)

	return &Renderer{
		Screen:      s,
		screenstyle: defStyle,
		cellstyle:   cellstyle,
		BorderStyle: borderstyle,
		Xoffset:     1,
		Yoffset:     1,
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
