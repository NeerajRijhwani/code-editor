package utils

import (
	"errors"
)

type GapBuffer struct {
	data     []rune
	gapStart int
	gapEnd   int
}

func Init_Gap_Buffer(text []rune, gapsize int) *GapBuffer {
	data := make([]rune, len(text)+gapsize)
	copy(data, text)
	return &GapBuffer{
		data:     data,
		gapStart: len(text),
		gapEnd:   len(text) + gapsize,
	}
}

func (g *GapBuffer) Len() int {
	return len(g.data) - g.GapSize()
}

func (g *GapBuffer) GapSize() int {
	return g.gapStart - g.gapEnd + 1
}

func (g *GapBuffer) Get(pos int) (rune, error) {
	if pos < 0 || pos >= g.GapSize() {
		return '_', errors.New("position Out of bounds")
	}
	if pos < g.gapStart {
		return g.data[pos], nil
	} else {
		return g.data[pos+g.gapStart], nil
	}
}

func (g *GapBuffer) MoveGap(pos int) error {
	if pos < 0 || pos > g.Len() {
		return errors.New("pos Out of bound")
	}
	if pos < g.gapStart {
		n := g.gapStart - pos
		copy(g.data[pos:g.gapStart], g.data[g.gapEnd-n:g.gapEnd])
		g.gapStart -= n
		g.gapEnd -= n
	} else {
		n := pos - g.gapStart
		copy(g.data[g.gapStart:g.gapStart+n], g.data[g.gapEnd:g.gapEnd+n])
		g.gapStart += n
		g.gapEnd += n
	}

	return nil
}

func (g *GapBuffer) Grow(n int) {
	data := make([]rune, n)
	copy(data[:g.gapStart], g.data[:g.gapStart])
	copy(data[g.gapEnd+n:], data[g.gapEnd:])
	g.gapEnd += n
	g.data = data
}

func (g *GapBuffer) Insert(pos int, ch rune) error {
	if err := g.MoveGap(pos); err != nil {
		return err
	}
	if g.GapSize() == 0 {
		g.Grow(12)
	}
	g.data[g.gapStart] = ch
	g.gapStart++
	return nil
}

func (g *GapBuffer) Delete(pos int) error {
	if err := g.MoveGap(pos); err != nil {
		return err
	}

	if g.gapEnd >= len(g.data) {
		return nil
	}

	g.gapEnd++
	return nil
}
func (g *GapBuffer) ToRunes() []rune {
	data := make([]rune, g.Len())

	copy(data[:g.gapStart], g.data[:g.gapStart])
	copy(data[g.gapStart:], g.data[g.gapEnd:])

	return data
}

func (g *GapBuffer) toString() string {
	return string(g.ToRunes())
}
