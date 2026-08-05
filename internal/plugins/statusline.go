package plugins

import (
	"fmt"

	"github.com/NeerajRijhwani/code-editor/internal/renderer"
	"github.com/gdamore/tcell/v3"
	"github.com/gdamore/tcell/v3/color"
)

const (
	NORMAL = 'n'
	INSERT = 'i'
	VISUAL = 'v'
)

type Status struct {
	height   int
	width    int
	startPos int
	mode     rune
	row      int
	col      int
	filename string
	saved    bool
	filetype string
	style    tcell.Style
}

func Init_StatusLine(height, width, start int, filename, filetype string) *Status {
	return &Status{
		height:   height,
		width:    width,
		startPos: start,
		mode:     NORMAL,
		row:      0,
		col:      0,
		filename: filename,
		saved:    true,
		filetype: filetype,
	}
}

func (s *Status) ChangeMode(mode rune) {
	switch mode {
	case NORMAL:
		s.mode = NORMAL
	case INSERT:
		s.mode = INSERT
	case VISUAL:
		s.mode = VISUAL
	default:
	}
}

func (s *Status) IsSaved() bool {
	return s.saved
}

func (s *Status) Unsaved() {
	s.saved = false
}
func (s *Status) Saved() {
	s.saved = true
}

func (s *Status) UpdateStatus(x, y int) {
	s.row, s.col = x, y
}

func (s *Status) RenderStatusLine(r *renderer.Renderer, x, y int) {
	s.UpdateStatus(x, y)
	mode := ""
	var modestyle tcell.Style
	switch s.mode {
	case NORMAL:
		mode = " NORMAL "
		modestyle = tcell.StyleDefault.Background(color.IndianRed).Foreground(color.Black.TrueColor())
	case INSERT:
		mode = " INSERT "
		modestyle = tcell.StyleDefault.Background(color.LightCyan).Foreground(color.Black.TrueColor())
	case VISUAL:
		mode = " VISUAL "
		modestyle = tcell.StyleDefault.Background(color.Orange).Foreground(color.Black.TrueColor())
	}
	r.DrawBox(0, s.startPos+s.height/2, s.width, s.startPos+s.height/2, tcell.StyleDefault.Background(color.SlateGrey))
	pos := 0
	for _, ch := range mode {
		r.Screen.SetContent(pos, s.startPos+(s.height/2), ch, nil, modestyle)
		pos++
	}
	filename := " " + s.filename + " "
	for _, ch := range filename {
		r.Screen.SetContent(pos, s.startPos+(s.height/2), ch, nil, s.style)
		pos++
	}

	if !s.saved {
		for _, ch := range "[+]" {
			r.Screen.SetContent(pos, s.startPos+(s.height/2), ch, nil, s.style)
			pos++
		}
	}
	cursortext := fmt.Sprintf(" %d:%d ", s.row, s.col)
	pos = s.width
	for i := len(cursortext) - 1; i >= 0; i-- {
		r.Screen.SetContent(pos, s.startPos+(s.height/2), rune(cursortext[i]), nil, modestyle)
		pos--
	}
	filetype := " " + s.filetype + " "
	for i := len(filetype) - 1; i >= 0; i-- {
		r.Screen.SetContent(pos, s.startPos+(s.height/2), rune(filetype[i]), nil, s.style)
		pos--
	}

}
