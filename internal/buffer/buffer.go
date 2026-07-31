package buffer

import (
	"errors"
	"log"
)

type Buffer struct {
	lines []string
}

func (b *Buffer) DeleteText(x1, y1, x2, y2 int) {
	if x1 > x2 || (x1 == x2 && y1 > y2) {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	if x1 == x2 {
		line := b.lines[x1]
		b.lines[x1] = line[:y1] + line[y2:]
		return
	}

	first := b.lines[x1][:y1]

	last := b.lines[x2][y2:]

	b.lines[x1] = first + last

	b.lines = append(b.lines[:x1+1], b.lines[x2+1:]...)

}

func (b *Buffer) GetSelectedText(x1, y1, x2, y2 int) string {
	res := ""

	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}
	if x1 == x2 {
		if y1 > y2 {
			y1, y2 = y2, y1
		}
		return b.lines[x1][y1:y2]
	}
	res += b.lines[x1][y1:] + string('\n')
	for i := x1 + 1; i < x2; i++ {
		res += b.lines[i] + string('\n')
	}
	res += b.lines[x2][:y2]
	return res
}

func InitBuffer(lines []string) *Buffer {
	if len(lines) == 0 {
		return &Buffer{
			lines: make([]string, 1),
		}
	}
	return &Buffer{
		lines: lines,
	}
}

func (b *Buffer) GetLine(i int) (string, error) {
	if i >= len(b.lines) {
		return "", errors.New("Invalid Line")
	}

	return b.lines[i], nil
}

func (b *Buffer) InsertEmptyLine(index int) {
	b.lines = append(b.lines, "")

	copy(b.lines[index+1:], b.lines[index:])

	b.lines[index] = ""
}

func (b *Buffer) InsertRune(row int, col int, ch rune) error {
	log.Printf("%d and %d", row, col)
	if len(b.lines) <= row || row < 0 {
		return errors.New("row exceeds buffer length")
	}

	if len(b.lines[row]) < col {
		return errors.New("col exceeds buffer length")
	}

	b.lines[row] = b.lines[row][:col] + string(ch) + b.lines[row][col:]

	return nil
}

func (b *Buffer) InsertText(text string, x, y int) error {
	if x >= len(b.lines) || x < 0 {
		return errors.New("Row Inavlid")
	}
	if y < 0 || y > len(b.lines[x]) {
		return errors.New("Col Inavlid")
	}
	b.lines[x] = b.lines[x][:y] + text + b.lines[x][y:]
	return nil
}

func (b *Buffer) DeleteRune(row int, col int) error {

	// log.Printf("%d and %d", row, col)
	// log.Printf("Current Buffer length: %d and col : %d", len(b.lines[row]),col)
	if len(b.lines) <= row {
		return errors.New("row exceeds buffer length")
	}

	if len(b.lines[row]) < col {
		return errors.New("col exceeds buffer length")
	}

	b.lines[row] = b.lines[row][:col] + b.lines[row][col+1:]

	return nil
}

func (b *Buffer) SplitLine(row int, col int) error {
	if len(b.lines) <= row {
		return errors.New("row exceeds buffer length")
	}

	if len(b.lines[row]) < col {
		return errors.New("col exceeds buffer length")
	}

	b.InsertEmptyLine(row + 1)
	if col != len(b.lines[row]) {

		b.lines[row+1] = b.lines[row][col:]
		b.lines[row] = b.lines[row][:col]

	}

	return nil
}

func (b *Buffer) MergeLine(row int) error {
	if len(b.lines) <= row {
		return errors.New("row exceeds buffer length")
	}

	if row != 0 {
		b.lines[row-1] = b.lines[row-1] + b.lines[row]
		b.lines = append(b.lines[:row], b.lines[row+1:]...)
	}
	return nil
}

func (b *Buffer) LineCount() int {
	return len(b.lines)
}

func (b *Buffer) LineLength(row int) (int, error) {

	if len(b.lines) <= row {
		return -1, errors.New("row exceeds buffer length")
	}

	return len(b.lines[row]), nil
}
