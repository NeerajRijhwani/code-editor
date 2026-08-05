package editor

import (
	"fmt"
	"log"
	"unicode/utf8"

	"github.com/NeerajRijhwani/code-editor/internal/buffer"
	"github.com/gdamore/tcell/v3"
)

func (e *Editor) HandleKey(ev *tcell.EventKey) {
	switch e.Mode {
	case 'n':
		switch ev.Key() {
		case tcell.KeyCtrlQ:
			e.running = false
		case tcell.KeyCtrlS:
			err := e.SaveFile()
			if err != nil {
				log.Printf("Unable to save file %v", err)
			}
		case tcell.KeyRight:
			e.RightKey()
		case tcell.KeyLeft:
			e.LeftKey()
		case tcell.KeyUp:
			e.UpKey()
		case tcell.KeyDown:
			e.DownKey()
		case tcell.KeyRune:
			ch, _ := utf8.DecodeRuneInString(ev.Str())
			switch ch {
			case 'i':
				log.Println("Insert Mode On")
				e.Mode = 'i'
			case 'v':
				log.Println("Visual Mode On")
				e.EnterVisualMode()
				e.Mode = 'v'
			case 'u':
				log.Println("Undo Pressed")
				e.History.Execute(nil, e.Buffer, 'u')
			case 'h':
				e.LeftKey()
			case 'j':
				e.DownKey()
			case 'k':
				e.UpKey()
			case 'l':
				e.RightKey()
			}
		case tcell.KeyCtrlR:
			e.History.Execute(nil, e.Buffer, 'r')
		}
	case 'i':
		log.Printf("key name pressed : %s ", ev.Name())
		switch ev.Key() {
		case tcell.KeyEscape:
			fallthrough
		case tcell.KeyCtrlQ:
			log.Println("Normal Mode On")
			e.Mode = 'n'
		case tcell.KeyCtrlS:
			err := e.SaveFile()
			if err != nil {
				log.Printf("Unable to save file %v", err)
			}
		case tcell.KeyRune:
			ch, _ := utf8.DecodeRuneInString(ev.Str())
			e.Insert(ch)
		case tcell.KeyEnter, tcell.KeyCtrlJ:
			e.Enter()
		case tcell.KeyBackspace:
			e.Backspace()
		case tcell.KeyDelete:
			e.Delete()
		case tcell.KeyRight:
			e.RightKey()
		case tcell.KeyLeft:
			e.LeftKey()
		case tcell.KeyUp:
			e.UpKey()
		case tcell.KeyDown:
			e.DownKey()
		}
	case 'v':
		x1, y1, x2, y2 := e.Select.Position()
		log.Printf("start %d, %d and end %d , %d", x1, y1, x2, y2)
		switch ev.Key() {
		case tcell.KeyEscape:
			log.Println("Normal Mode On")
			e.Select.Reset()
			e.Mode = 'n'
		case tcell.KeyRight:
			e.RightKey()
		case tcell.KeyLeft:
			e.LeftKey()
		case tcell.KeyUp:
			e.UpKey()
		case tcell.KeyDown, tcell.KeyEnter:
			e.DownKey()
		case tcell.KeyRune:
			ch, _ := utf8.DecodeRuneInString(ev.Str())
			switch ch {
			case 'y':
				e.Yank()
			case 'd':
				e.Cut()
			}
		}
	}
}

func (e *Editor) Yank() {
	text := e.Buffer.GetSelectedText(e.Select.Position())
	e.Renderer.SetClipboard(text)
	e.Select.Reset()
	e.Mode = 'n'
}

func (e *Editor) Cut() {
	x1, y1, x2, y2 := e.Select.Position()
	e.Buffer.DeleteText(x1, y1, x2, y2)
	if x1 > x2 {
		e.Cursor.SetCursor(x2, y2)
	} else {
		e.Cursor.SetCursor(x1, y1)
	}
	e.Select.Reset()
	e.Mode = 'n'
}

func (e *Editor) EnterVisualMode() {
	x, y := e.Cursor.Position()
	e.Select.Active = true
	e.Select.SetStartCoord(x, y)
	e.Select.SetEndCoord(x, y)
}

func (e *Editor) Insert(key rune) {
	x, y := e.Cursor.Position()
	e.Buffer.InsertRune(x, y, key)
	cmd := buffer.Init_InsertTextCommand(x, y, x, y+1, string(key))
	e.History.Execute(cmd, e.Buffer, 'i')
	e.Cursor.MoveRight()
}
func (e *Editor) Backspace() {
	log.Println("Backspace Pressed")
	x, y := e.Cursor.Position()
	if x == 0 && y == 0 {
	} else if y != 0 {
		ch, err := e.Buffer.DeleteRune(x, y-1)
		if err != nil {
			log.Printf("Backspace error: %v", err)
		}
		cmd := buffer.Init_DeleteTextCommand(x, y-1, x, y, ch, 0)
		log.Printf("Command init : %v ", cmd)
		e.History.Execute(cmd, e.Buffer, 'd')

		e.Cursor.MoveLeft()
	} else {
		length, _ := e.Buffer.LineLength(x - 1)
		err := e.Buffer.MergeLine(x)
		if err != nil {
			log.Printf("Merge Line error: %v", err)
		}
		cmd := buffer.Init_MergeLineCommand(x, 0, length)
		e.History.Execute(cmd, e.Buffer, 'm')
		e.Cursor.SetCursor(x-1, length)
	}
}

func (e *Editor) Delete() {
	log.Println("Delete Pressed")
	x, y := e.Cursor.Position()
	length, _ := e.Buffer.LineLength(x)
	if x == e.Buffer.LineCount() && y == length {
		// do nothing
	} else if y != length {
		ch, err := e.Buffer.DeleteRune(x, y)
		if err != nil {
			log.Printf("Delete error: %v", err)
		}
		cmd := buffer.Init_DeleteTextCommand(x, y, x, y+1, ch, 1)
		e.History.Execute(cmd, e.Buffer, 'd')

	} else {
		err := e.Buffer.MergeLine(x + 1)
		if err != nil {
			log.Printf("Merge Line error: %v", err)
		}
		cmd := buffer.Init_MergeLineCommand(x+1, 0, length)
		e.History.Execute(cmd, e.Buffer, 'm')

	}
}

func (e *Editor) Enter() {
	log.Println("Enter Pressed")
	x, y := e.Cursor.Position()
	err := e.Buffer.SplitLine(x, y)
	if err != nil {
		fmt.Printf("%v", err)
	}
	cmd := buffer.Init_SplitLineCommand(x, y)
	e.History.Execute(cmd, e.Buffer, 's')
	if x == 30 {
		e.Renderer.IncreaseFirstLine()
	}
	e.Cursor.SetCursor(x+1, 0)
	screenRow := (x + 1) - e.Renderer.FirstLine
	if screenRow >= 30 {
		e.Renderer.FirstLine++
	}
}

func (e *Editor) LeftKey() {
	log.Println("LeftKey Pressed")
	_, y := e.Cursor.Position()
	if y != 0 {
		e.Cursor.MoveLeft()
	}
	if e.Select.Active {
		currx, curry := e.Cursor.Position()
		e.Select.SetEndCoord(currx, curry)
	}
}

func (e *Editor) RightKey() {
	log.Println("RightKey Pressed")
	x, y := e.Cursor.Position()
	length, _ := e.Buffer.LineLength(x)
	if y != length {
		e.Cursor.MoveRight()
	}
	if e.Select.Active {
		currx, curry := e.Cursor.Position()
		e.Select.SetEndCoord(currx, curry)
	}
}

func (e *Editor) UpKey() {
	log.Println("UpKey Pressed")

	x, y := e.Cursor.Position()
	if x == 0 {
		return
	}
	prevLen, _ := e.Buffer.LineLength(x - 1)

	if y > prevLen {
		e.Cursor.SetCursor(x, prevLen)
	}
	e.Cursor.MoveUp()
	screenRow := (x - 1) - e.Renderer.FirstLine

	if screenRow < 0 {
		e.Renderer.FirstLine--
	}
	if e.Select.Active {
		currx, curry := e.Cursor.Position()
		e.Select.SetEndCoord(currx, curry)
	}

}

func (e *Editor) DownKey() {
	log.Println("DownKey Pressed")

	x, y := e.Cursor.Position()

	if x == e.Buffer.LineCount()-1 {
		return
	}

	nextLen, _ := e.Buffer.LineLength(x + 1)
	if y > nextLen {
		e.Cursor.SetCursor(x, nextLen)
	}
	e.Cursor.MoveDown()

	screenRow := (x + 1) - e.Renderer.FirstLine

	if screenRow >= 30 {
		e.Renderer.FirstLine++
	}
	if e.Select.Active {
		currx, curry := e.Cursor.Position()
		e.Select.SetEndCoord(currx, curry)
	}
}
