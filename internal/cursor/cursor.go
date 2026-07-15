package cursor

type Cursor struct {
	x int
	y int
}

func initCursor() *Cursor {
	return &Cursor{
		x: 0,
		y: 0,
	}
}

func (c *Cursor) MoveRight() {
	c.x++
}

func (c *Cursor) MoveLeft() {
	c.x--
}

func (c *Cursor) MoveUp() {
	c.y--
}

func (c *Cursor) MoveDown() {
	c.y++
}

func (c *Cursor) Position() (int, int) {
	return c.x, c.y
}

func (c *Cursor) SetCursor(x, y int) {
	c.x = x
	c.y = y
}
