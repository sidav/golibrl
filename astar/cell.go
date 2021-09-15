package astar

type Cell struct {
	X, Y            int
	g, h            int
	parent          *Cell
	Child           *Cell
	numChilds 		int
}

func (c *Cell) getF() int {
	return c.g + c.h
}

func (c *Cell) GetCoords() (int, int) {
	return c.X, c.Y
}

func (c *Cell) setG(inc int) {
	if c.parent != nil {
		c.g = c.parent.g + inc
	}
}

func (c *Cell) GetNextStepVector() (int, int) {
	var x, y int
	if c.Child != nil {
		x = c.Child.X - c.X
		y = c.Child.Y - c.Y
	}
	return x, y
}

//func (c *Cell) getPathToCell() *[]*Cell {
//	path := make([]*Cell, 0)
//	curcell := c
//	for curcell != nil {
//		path = append(path, curcell)
//		curcell = curcell.parent
//	}
//	return &path
//}

func (c *Cell) setChildsForPath() {
	// path := make([]*Cell, 0)
	curcell := c
	c.numChilds = 0
	for curcell.parent != nil {
		// path = append(path, curcell)
		curcell.parent.Child = curcell
		curcell.parent.numChilds = curcell.numChilds+1
		curcell = curcell.parent
	}
	return
}

func (c *Cell) GetTotalPathLength() int {
	return c.numChilds
}
