package main

import (
	"fmt"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
	"log"
)

type Editor struct {
	lines []string

	cursorX int
	cursorY int
	scrollX int
	scrollY int
}

func initEditor() *Editor {

	return &Editor{
		lines:   make([]string, 1),
		cursorX: 1,
		cursorY: 1,
		scrollX: 1,
		scrollY: 1,
	}
}

func main() {
	fmt.Println("hello")
	initScreen()

}

func drawBox(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	if y2 < y1 {
		y1, y2 = y2, y1
	}
	if x2 < x1 {
		x1, x2 = x2, x1
	}

	// Fill background
	for row := y1; row <= y2; row++ {
		for col := x1; col <= x2; col++ {
			s.Put(col, row, " ", style)
		}
	}

	// Draw borders
	for col := x1; col <= x2; col++ {
		s.Put(col, y1, string(tcell.RuneHLine), style)
		s.Put(col, y2, string(tcell.RuneHLine), style)
	}
	for row := y1 + 1; row < y2; row++ {
		s.Put(x1, row, string(tcell.RuneVLine), style)
		s.Put(x2, row, string(tcell.RuneVLine), style)
	}

	// Only draw corners if necessary
	if y1 != y2 && x1 != x2 {
		s.Put(x1, y1, string(tcell.RuneULCorner), style)
		s.Put(x2, y1, string(tcell.RuneURCorner), style)
		s.Put(x1, y2, string(tcell.RuneLLCorner), style)
		s.Put(x2, y2, string(tcell.RuneLRCorner), style)
	}

	drawText(s, x1+1, y1+1, x2-1, y2-1, style, text)
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	var width int
	for text != "" {
		text, width = s.Put(col, row, text, style)
		col += width
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
		if width == 0 {
			// incomplete grapheme at end of string
			break
		}
	}
}

func renderText(e *Editor, s tcell.Screen, style tcell.Style) {
	value := e.lines[len(e.lines)-1]
	drawBox(s, e.cursorX, e.cursorY, e.cursorX+len(value), e.cursorY, style, value)

	e.cursorX++
	if e.cursorX > 10 {
		e.cursorY++
		e.cursorX = 0
	}
}

func updateLine(e *Editor, key string) {
	e.lines[e.cursorY-1] += key
	fmt.Println(key)
}

func initScreen() {
	boxStyle := tcell.StyleDefault.Foreground(color.White).Background(color.Black)
	defStyle := tcell.StyleDefault.Background(color.Reset).Foreground(color.Reset)

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
		log.Fatalf("%+v", err)
	}

	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	// Set default text style

	s.SetStyle(defStyle)
	s.EnableMouse()
	s.EnablePaste()
	// Clear screen
	s.Clear()

	editor := initEditor()

	defer quit(s)

	// ox, oy := -1, -1

	for {

		s.Show()
		ev := <-s.EventQ()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				return
			} else if ev.Key() == tcell.KeyCtrlL {
				s.Sync()
			} else {
				updateLine(editor, ev.Name())
				renderText(editor, s, boxStyle)
			}
		}
	}
}

func quit(s tcell.Screen) {
	// You have to catch panics in a defer, clean up, and
	// re-raise them - otherwise your application can
	// die without leaving any diagnostic trace.
	maybePanic := recover()
	s.Fini()
	if maybePanic != nil {
		panic(maybePanic)
	}
}
