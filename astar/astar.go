package astar

const (
	DIAGONAL_COST = 14
	STRAIGHT_COST = 10
	HEURISTIC_MULTIPLIER = 10
	MAX_PATHFINDING_STEPS = 175 // Increase in case of stupid pathfinding. Decrease in case of lag.
)

type Cell struct {
	X, Y            int
	g, h            int
	costToMoveThere int
	parent          *Cell
	Child           *Cell
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func heuristicCost(fromx, fromy, tox, toy int, diagonalsAllowed bool) int {
	if diagonalsAllowed {
		return (fromx-tox)*(fromx-tox)+(fromy-toy)*(fromy-toy)
	}
	return HEURISTIC_MULTIPLIER * (abs(tox-fromx) + abs(toy-fromy))
}

func getIndexOfCellWithLowestF(openList []*Cell) int {
	cheapestCellIndex := 0
	for i, c := range openList {
		if c.getF() < openList[cheapestCellIndex].getF() {
			cheapestCellIndex = i
		}
	}
	return cheapestCellIndex
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
	for curcell.parent != nil {
		// path = append(path, curcell)
		curcell.parent.Child = curcell
		curcell = curcell.parent
	}
	return
}

func FindPath(costMap *[][]int, fromx, fromy, tox, toy int, diagonalMoveAllowed, forceGetPath bool) *Cell {
	openList := make([]*Cell, 0)
	closedList := make([]*Cell, 0)
	var currentCell *Cell
	total_steps := 0
	targetReached := false

	// step 1
	origin := &Cell{X: fromx, Y: fromy, costToMoveThere: 0, h: heuristicCost(fromx, fromy, tox, toy, diagonalMoveAllowed)}
	openList = append(openList, origin)
	// step 2
	for !targetReached {
		// sub-step 2a:
		currentCellIndex := getIndexOfCellWithLowestF(openList)
		currentCell = openList[currentCellIndex]
		// sub-step 2b:
		closedList = append(closedList, currentCell)
		openList = append(openList[:currentCellIndex], openList[currentCellIndex+1:]...) // this friggin' magic removes currentCellIndex'th element from openList
		//sub-step 2c:
		analyzeNeighbors(currentCell, &openList, &closedList, costMap, tox, toy, diagonalMoveAllowed)
		//sub-step 2d:
		total_steps += 1
		targetInOpenList := getCellWithCoordsFromList(&openList, tox, toy)
		if targetInOpenList != nil {
			currentCell = targetInOpenList
			currentCell.setChildsForPath()
			return origin
		}
		if len(openList) == 0 || total_steps > MAX_PATHFINDING_STEPS {
			if forceGetPath { // makes the routine always return path to the closest possible cell to (tox, toy) even if the precise path does not exist.
				currentCell = getCellWithLowestHeuristicFromList(&closedList)
				currentCell.setChildsForPath()
				return origin
			} else {
				return nil
			}
		}
	}
	return nil
}

func analyzeNeighbors(curCell *Cell, openlist *[]*Cell, closedlist *[]*Cell, costMap *[][]int, targetX, targetY int, diagAllowed bool) {
	cost := 0
	cx, cy := curCell.X, curCell.Y
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			x, y := cx+i, cy+j
			if areCoordsValidForCostMap(x, y, costMap) {
				// if (x != targetX || y != targetY) &&
				if (*costMap)[x][y] == -1 || getCellWithCoordsFromList(closedlist, x, y) != nil { // Cell is impassable or is in closed list
					continue // ignore it
				}
				// TODO: add a flag for skipping diagonally lying cells
				// TODO: add actual "cost to move there" from costMap
				if (i * j) != 0 { // the Cell under consideration is lying diagonally
					cost = DIAGONAL_COST
				} else {
					cost = STRAIGHT_COST
				}
				curNeighbor := getCellWithCoordsFromList(openlist, x, y)
				if curNeighbor != nil {
					if curNeighbor.g > curCell.g+cost {
						curNeighbor.parent = curCell
						curNeighbor.setG(cost)
					}
				} else {
					curNeighbor = &Cell{X: x, Y: y, parent: curCell, h: heuristicCost(x, y, targetX, targetY, diagAllowed)}
					curNeighbor.setG(cost)
					*openlist = append(*openlist, curNeighbor)
				}
			}
		}
	}
}

func getCellWithCoordsFromList(list *[]*Cell, x, y int) *Cell {
	for _, c := range *list {
		if c.X == x && c.Y == y {
			return c
		}
	}
	return nil
}

func getCellWithLowestHeuristicFromList(list *[]*Cell) *Cell {
	lowest := (*list)[0]
	for _, c := range *list {
		if c.h < lowest.h {
			lowest = c
		}
	}
	return lowest
}

func areCoordsValidForCostMap(x, y int, costMap *[][]int) bool {
	return x >= 0 && y >= 0 && (x < len(*costMap)) && (y < len((*costMap)[0]))
}

// АЛГОРИТМ А*:
// Замечание: F = g+h, где g - цена пути ИЗ стартовой точки, h - эвристическая оценка пути ДО цели.
// По "методу Манхэттена": h = 10 * (_abs(targetX - startX) + _abs(targetY-startY))
//  1) Добавляем стартовую клетку в открытый список.
//  2) Повторяем следующее:
//  a) Ищем в открытом списке клетку с наименьшей стоимостью F. Делаем ее текущей клеткой.
//  b) Помещаем ее в закрытый список. (И удаляем с открытого)
//  c) Для каждой из соседних 8-ми клеток ...
//  	Если клетка непроходимая или она находится в закрытом списке, игнорируем ее. В противном случае делаем следующее.
//  	Если клетка еще не в открытом списке, то добавляем ее туда. Делаем текущую клетку родительской для это клетки. Расчитываем стоимости F, g и h клетки.
//  	Если клетка уже в открытом списке, то проверяем, не дешевле ли будет путь через эту клетку. Для сравнения используем стоимость g.
//  	Более низкая стоимость g указывает на то, что путь будет дешевле. Эсли это так, то меняем родителя клетки на текущую клетку и пересчитываем для нее стоимости g и F. Если вы сортируете открытый список по стоимости F, то вам надо отсортировать свесь список в соответствии с изменениями.
//  d) Останавливаемся если:
//  	Добавили целевую клетку в открытый список, в этом случае путь найден.
// 		Или открытый список пуст и мы не дошли до целевой клетки. В этом случае путь отсутствует.
//  3) Сохраняем путь. Двигаясь назад от целевой точки, проходя от каждой точки к ее родителю до тех пор, пока не дойдем до стартовой точки. Это и будет наш путь.
