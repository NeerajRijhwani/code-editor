package editor

import (
	"fmt"
	"log"
	// "log"
	"unicode/utf8"

	"github.com/gdamore/tcell/v3"
)

func (e *Editor) HandleKey(ev *tcell.EventKey) {
	switch ev.Key() {
	case tcell.KeyEscape:
		fallthrough
	case tcell.KeyCtrlQ:
		e.running = false
	case tcell.KeyRune:
		ch, _ := utf8.DecodeRuneInString(ev.Str())
		e.Insert(ch)
	case tcell.KeyEnter:
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
}

func (e *Editor) Insert(key rune) {
	log.Printf("Key Pressed : %c", key)
	x, y := e.Cursor.Position()
	e.Buffer.InsertRune(x, y, key)
	e.Cursor.MoveRight()
}

func (e *Editor) Backspace() {
	log.Println("Backspace Pressed")
	x, y := e.Cursor.Position()
	if x == 0 && y == 0 {
	} else if y != 0 {
		err := e.Buffer.DeleteRune(x, y-1)
		if err != nil {
			log.Printf("Backspace error: %v", err)
		}
		e.Cursor.MoveLeft()
	} else {
		length, _ := e.Buffer.LineLength(x - 1)
		err := e.Buffer.MergeLine(x)
		if err != nil {
			log.Printf("Merge Line error: %v", err)
		}
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
		err := e.Buffer.DeleteRune(x, y)
		if err != nil {
			log.Printf("Delete error: %v", err)
		}
	} else {
		err := e.Buffer.MergeLine(x + 1)
		if err != nil {
			log.Printf("Merge Line error: %v", err)
		}
	}
}

func (e *Editor) Enter() {
	log.Println("Enter Pressed")
	x, y := e.Cursor.Position()
	err := e.Buffer.SplitLine(x, y)
	if err != nil {
		fmt.Printf("%v", err)
	}
	e.Cursor.SetCursor(x+1, 0)
}

func (e *Editor) LeftKey() {
	log.Println("LeftKey Pressed")
	_, y := e.Cursor.Position()
	if y != 0 {
		e.Cursor.MoveLeft()
	}
}

func (e *Editor) RightKey() {
	log.Println("RightKey Pressed")
	x, y := e.Cursor.Position()
	length, _ := e.Buffer.LineLength(x)
	if y != length {
		e.Cursor.MoveRight()
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

}

func (e *Editor) DownKey() {
	log.Println("DownKey Pressed")

	x, y := e.Cursor.Position()
	count := e.Buffer.LineCount()

	if x >= count-1 {
		return
	}

	nextLen, _ := e.Buffer.LineLength(x + 1)
	if y > nextLen {
		e.Cursor.SetCursor(x, nextLen)
	}
	e.Cursor.MoveDown()
}
