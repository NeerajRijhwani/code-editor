package cursor

// import "fmt"

type Cursor struct {
	x int
	y int
}

func InitCursor() *Cursor {
	return &Cursor{
		x: 0,
		y: 0,
	}
}

func (c *Cursor) MoveRight() {
	c.y++
}

func (c *Cursor) MoveLeft() {
	c.y--
}

func (c *Cursor) MoveUp() {
	c.x--
}

func (c *Cursor) MoveDown() {
	c.x++
}

func (c *Cursor) Position() (int, int) {
	// fmt.Printf("X: %d  Y: %d \n", c.x, c.y)
	return c.x, c.y
}

func (c *Cursor) SetCursor(x, y int) {
	c.x = max(0, x)
	c.y = max(0, y)
}
