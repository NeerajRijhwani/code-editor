package cursor

type Selection struct {
	Active bool

	StartX int
	StartY int

	EndX int
	EndY int
}

func InitSelect() *Selection {
	return &Selection{
		Active: false,
		StartX: -1,
		StartY: -1,
		EndX:   -1,
		EndY:   -1,
	}
}

func (s *Selection) Reset() {
	s.Active = false
	s.StartX = -1
	s.StartY = -1
	s.EndX = -1
	s.EndY = -1
}

func (s *Selection) Position() (int, int, int, int) {
	return s.StartX, s.StartY, s.EndX, s.EndY
}

func (s *Selection) CheckWithinSelect(x int, y int) bool {
	sX, sY, eX, eY := s.Position()
	if sX > eX {
		sX, eX = eX, sX
		sY, eY = eY, sY

	}
	if x > sX && x < eX {
		return true
	} else if x == sX && x == eX {
		if y >= sY && y <= eY {
			return true
		} else {
			return false
		}

	} else if x == sX {
		if y >= sY {
			return true
		} else {
			return false
		}
	} else if x == eX {
		if y <= eY {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (s *Selection) SetStartCoord(x, y int) {
	s.StartX = x
	s.StartY = y
}
func (s *Selection) SetEndCoord(x, y int) {
	s.EndX = x
	s.EndY = y
}
