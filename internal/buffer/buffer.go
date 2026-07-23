package buffer

import (
	"errors"
	"log"
	// "fmt"
)

type Buffer struct {
	lines []string
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
