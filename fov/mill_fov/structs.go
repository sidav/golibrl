package mill_fov

type slope struct {
	Y, X int
}

func (s *slope) Greater(y, x int) bool {
	return s.Y*x > s.X*y
}

func (s *slope) GreaterOrEqual(y, x int) bool {
	return s.Y*x >= s.X*y
}

func (s *slope) Less(y, x int) bool {
	return s.Y*x < s.X*y
}