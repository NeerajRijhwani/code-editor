package buffer

import (
	"errors"
	"log"
	"strings"

	"github.com/NeerajRijhwani/code-editor/internal/utils"
)

const (
	GAP_SIZE = 10
)

type Buffer struct {
	lines []*utils.GapBuffer
}

func (b *Buffer) DeleteText(x1, y1, x2, y2 int) {
	if x1 > x2 || (x1 == x2 && y1 > y2) {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	if x1 == x2 {
		b.lines[x1].DeleteText(y1, y2-y1)
		return
	}

	text := b.lines[x2].ToString()[y2:]
	length := b.lines[x1].Len()
	b.lines[x1].DeleteText(y1, length-1-y1)
	b.lines[x1].InsertText(y2, text)

	b.lines = append(b.lines[:x1+1], b.lines[x2+1:]...)
}

func (b *Buffer) GetSelectedText(x1, y1, x2, y2 int) string {

	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	if x1 == x2 {
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		return b.lines[x1].ToString()[y1:y2]
	}
	var sb strings.Builder
	sb.WriteString(b.lines[x1].ToString())
	sb.WriteString(string('\n'))
	for i := x1 + 1; i < x2; i++ {
		sb.WriteString(b.lines[i].ToString())
		sb.WriteString(string('\n'))
	}
	sb.WriteString(b.lines[x2].ToString()[:y2])
	return sb.String()
}

func InitBuffer(lines []string) *Buffer {
	if len(lines) == 0 {
		line := utils.Init_Gap_Buffer(make([]rune, 0), GAP_SIZE)

		buf := Buffer{
			lines: make([]*utils.GapBuffer, 1),
		}
		buf.lines[0] = line
		return &buf
	}

	buf := Buffer{
		lines: make([]*utils.GapBuffer, len(lines)),
	}

	for i := range len(lines) {
		line := utils.Init_Gap_Buffer([]rune(lines[i]), GAP_SIZE)
		buf.lines[i] = line
	}

	return &buf
}

func (b *Buffer) GetBuffer() []byte {
	var sb strings.Builder
	for i := range len(b.lines) {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(b.lines[i].ToString())
	}
	buf := []byte(sb.String())
	return buf
}

func (b *Buffer) GetLine(i int) (string, error) {
	if i >= len(b.lines) {
		return "", errors.New("Invalid Line")
	}
	line := b.lines[i].ToString()

	return line, nil
}

func (b *Buffer) InsertEmptyLine(index int) {
	empty := make([]rune, 0)
	line := utils.Init_Gap_Buffer(empty, GAP_SIZE)
	b.lines = append(b.lines, line)
	copy(b.lines[index+1:], b.lines[index:])
	b.lines[index] = line
}

func (b *Buffer) InsertRune(row int, col int, ch rune) error {
	log.Printf("%d and %d", row, col)
	if len(b.lines) <= row || row < 0 {
		return errors.New("row exceeds buffer length")
	}
	length := b.lines[row].Len()
	if length < col {
		return errors.New("col exceeds buffer length")
	}

	if err := b.lines[row].Insert(col, ch); err != nil {
		return err
	}

	return nil
}

func (b *Buffer) InsertText(text string, x, y int) error {
	log.Printf("InsertText called %s with x %d and y %d ", text, x, y)
	log.Printf("col : %d", b.lines[x].Len())
	if x >= len(b.lines) || x < 0 {
		return errors.New("Row Inavlid")
	}
	if y < 0 || y > b.lines[x].Len() {
		return errors.New("Col Inavlid")
	}
	if err := b.lines[x].InsertText(y, text); err != nil {
		return err
	}
	return nil
}

func (b *Buffer) DeleteRune(row int, col int) (string, error) {

	// log.Printf("%d and %d", row, col)
	// log.Printf("Current Buffer length: %d and col : %d", len(b.lines[row]),col)
	if len(b.lines) <= row {
		return "", errors.New("row exceeds buffer length")
	}
	length := b.lines[row].Len()
	if length < col {
		return "", errors.New("col exceeds buffer length")
	}
	ch, err := b.lines[row].Delete(col)
	if err != nil {
		return "", err
	}
	return string(ch), nil
}

func (b *Buffer) SplitLine(row int, col int) error {
	if len(b.lines) <= row {
		return errors.New("row exceeds buffer length")
	}

	if b.lines[row].Len() < col {
		return errors.New("col exceeds buffer length")
	}

	b.InsertEmptyLine(row + 1)
	if col != b.lines[row].Len() {
		text := b.lines[row].ToString()[col:]
		b.lines[row+1].InsertText(0, text)
		b.lines[row].DeleteText(col, len(text))
	}

	return nil
}

func (b *Buffer) MergeLine(row int) error {
	if len(b.lines) <= row {
		return errors.New("row exceeds buffer length")
	}

	if row != 0 {
		line := b.lines[row].ToString()
		pos := b.lines[row-1].Len()
		b.lines[row-1].InsertText(pos, line)
		b.lines = append(b.lines[:row], b.lines[row+1:]...)
	}
	return nil
}

func (b *Buffer) LineCount() int {
	return len(b.lines)
}

func (b *Buffer) LineLength(row int) (int, error) {

	if row < 0 || len(b.lines) <= row {
		return -1, errors.New("row Invalid")
	}

	length := b.lines[row].Len()
	return length, nil
}
