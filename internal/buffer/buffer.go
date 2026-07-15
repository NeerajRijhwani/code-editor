package buffer

import (
	"errors"
)

type Buffer struct {
	lines []string
}

func initBuffer() *Buffer {
	return &Buffer{
		lines: make([]string, 1),
	}
}

func (b *Buffer) InsertEmptyLine(index int) {
	b.lines = append(b.lines, "")

	copy(b.lines[index+1:], b.lines[index:])

	b.lines[index] = ""
}

func (b *Buffer) InsertRune(col, row int, ch rune) error {
	if len(b.lines) <= row {
		return errors.New("row exceeds buffer length")
	}

	if len(b.lines[row]) > col {
		return errors.New("col exceeds buffer length")
	}

	b.lines[row] = b.lines[row][:col] + string(ch) + b.lines[row][col:]

	return nil
}

func (b *Buffer) DeleteRune(col, row int) error {
	if len(b.lines) <= row {
		return errors.New("row exceeds buffer length")
	}

	if len(b.lines[row]) > col {
		return errors.New("col exceeds buffer length")
	}

	b.lines[row] = b.lines[row][:col] + b.lines[row][col+1:]

	return nil
}

func (b *Buffer) SplitLine(col, row int) error {
	if len(b.lines) <= row {
		return errors.New("row exceeds buffer length")
	}

	if len(b.lines[row]) > col {
		return errors.New("col exceeds buffer length")
	}

	b.InsertEmptyLine(row + 1)
	if col != len(b.lines[row]) {

		b.lines[row+1] = b.lines[row][col:]
		b.lines[row] = b.lines[row][:col]

	}

	return nil
}

func (b *Buffer) MergeLine(col, row int) error {
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

func (b *Buffer) LineLength(row int) (error, int) {

	if len(b.lines) <= row {
		return errors.New("row exceeds buffer length"), -1
	}

	return nil, len(b.lines[row])
}
