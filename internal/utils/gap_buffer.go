package utils

import (
	"errors"
	"log"
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
	return g.gapEnd - g.gapStart
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
	log.Printf(
		"MoveGap: pos=%d start=%d end=%d dataLen=%d len=%d\n",
		pos,
		g.gapStart,
		g.gapEnd,
		len(g.data),
		g.Len(),
	)
	if pos < g.gapStart {
		n := g.gapStart - pos
		copy(g.data[g.gapEnd-n:g.gapEnd], g.data[pos:g.gapStart])
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
	data := make([]rune, n+len(g.data))
	copy(data[:g.gapStart], g.data[:g.gapStart])
	copy(data[g.gapEnd+n:], g.data[g.gapEnd:])
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

func (g *GapBuffer) InsertText(pos int, text string) error {
	if err := g.MoveGap(pos); err != nil {
		return err
	}
	textrunes := []rune(text)
	length := len(textrunes)

	if g.GapSize() < length {
		g.Grow(length - g.GapSize() + 12)
	}

	copy(g.data[g.gapStart:], textrunes)

	g.gapStart += length
	return nil
}

func (g *GapBuffer) Delete(pos int) (rune, error) {
	if err := g.MoveGap(pos); err != nil {
		return ' ', err
	}

	if g.gapEnd >= len(g.data) {
		return ' ', nil
	}
	ch := g.data[pos]
	g.gapEnd++
	return ch, nil
}

func (g *GapBuffer) DeleteText(pos int, n int) error {
	if pos < 0 || n < 0 || pos+n > g.Len() {
		return errors.New("delete range out of bounds")
	}

	if err := g.MoveGap(pos); err != nil {
		return err
	}
	g.gapEnd += n
	return nil
}

func (g *GapBuffer) ToRunes() []rune {
	length := g.Len()
	data := make([]rune, length)
	copy(data[:g.gapStart], g.data[:g.gapStart])
	copy(data[g.gapStart:], g.data[g.gapEnd:])

	return data
}

func (g *GapBuffer) ToString() string {
	return string(g.ToRunes())
}
