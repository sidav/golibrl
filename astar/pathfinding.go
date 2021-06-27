package astar

const (
	DIAGONAL_COST             = 14
	STRAIGHT_COST             = 10
	HEURISTIC_MULTIPLIER      = 10
	DEFAULT_PATHFINDING_STEPS = 175 // Increase in case of stupid pathfinding. Decrease in case of lag.
)

func heuristicCost(fromx, fromy, tox, toy int, diagonalsAllowed bool) int {
	if diagonalsAllowed {
		return (fromx-tox)*(fromx-tox) + (fromy-toy)*(fromy-toy)
	}
	return HEURISTIC_MULTIPLIER * (abs(tox-fromx) + abs(toy-fromy))
}

type AStarPathfinder struct {
	DiagonalMoveAllowed, ForceGetPath, ForceIncludeFinish, AutoAdjustDefaultMaxSteps bool
	// private
	targetCell *Cell
	toX, toY int
}

func (aspf *AStarPathfinder) getIndexOfCellWithLowestF(openList []*Cell) int {
	cheapestCellIndex := 0
	for i, c := range openList {
		if c.getF() < openList[cheapestCellIndex].getF() {
			cheapestCellIndex = i
		}
	}
	return cheapestCellIndex
}

func (aspf *AStarPathfinder) FindPath(costMap *[][]int, fromx, fromy, tox, toy int) *Cell {
	aspf.toX = tox
	aspf.toY = toy
	aspf.targetCell = nil
	openList := make([]*Cell, 0)
	closedList := make([]*Cell, 0)
	var currentCell *Cell
	total_steps := 0
	targetReached := false
	totalCells := len(*costMap) * len((*costMap)[0])

	maxSearchDepth := DEFAULT_PATHFINDING_STEPS
	if totalCells > DEFAULT_PATHFINDING_STEPS && aspf.AutoAdjustDefaultMaxSteps {
		maxSearchDepth = totalCells
	}

	// step 1
	origin := &Cell{X: fromx, Y: fromy, h: heuristicCost(fromx, fromy, tox, toy, aspf.DiagonalMoveAllowed)}
	openList = append(openList, origin)
	// step 2
	for !targetReached {
		// sub-step 2a:
		currentCellIndex := aspf.getIndexOfCellWithLowestF(openList)
		currentCell = openList[currentCellIndex]
		// sub-step 2b:
		closedList = append(closedList, currentCell)
		openList = append(openList[:currentCellIndex], openList[currentCellIndex+1:]...) // this friggin' magic removes currentCellIndex'th element from openList
		//sub-step 2c:
		aspf.analyzeNeighbors(currentCell, &openList, &closedList, costMap, tox, toy)
		//sub-step 2d:
		total_steps += 1
		if aspf.targetCell != nil {
			aspf.targetCell.setChildsForPath()
			return origin
		}
		if len(openList) == 0 || total_steps > maxSearchDepth {
			if aspf.ForceGetPath { // makes the routine always return path to the closest possible cell to (tox, toy) even if the precise path does not exist.
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

func (aspf *AStarPathfinder) analyzeNeighbors(curCell *Cell, openlist *[]*Cell, closedlist *[]*Cell, costMap *[][]int, targetX, targetY int) {
	cost := 0
	cx, cy := curCell.X, curCell.Y
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if (i == 0 && j == 0) || (!aspf.DiagonalMoveAllowed && i != 0 && j != 0) {
				continue
			}
			x, y := cx+i, cy+j
			if areCoordsValidForCostMap(x, y, costMap) {
				// if (x != targetX || y != targetY) &&
				if (*costMap)[x][y] < 0 || getCellWithCoordsFromList(closedlist, x, y) != nil { // Cell is impassable or is in closed list
					if !(aspf.ForceIncludeFinish && x == targetX && y == targetY) { // if ForceIncludeFinish is true, then we won't ignore finish cell whether it is passable or whatever.
						continue // ignore it
					}
				}
				// TODO: add actual "cost to move there" from costMap
				if (i * j) != 0 { // the Cell under consideration is lying diagonally
					cost = DIAGONAL_COST * (*costMap)[x][y]
				} else {
					cost = STRAIGHT_COST * (*costMap)[x][y]
				}
				curNeighbor := getCellWithCoordsFromList(openlist, x, y)
				if curNeighbor != nil {
					if curNeighbor.g > curCell.g+cost {
						curNeighbor.parent = curCell
						curNeighbor.setG(cost)
					}
				} else {
					curNeighbor = &Cell{X: x, Y: y, parent: curCell, h: heuristicCost(x, y, targetX, targetY, aspf.DiagonalMoveAllowed)}
					if x == aspf.toX && y == aspf.toY {
						aspf.targetCell = curNeighbor
						return 
					}
					curNeighbor.setG(cost)
					*openlist = append(*openlist, curNeighbor)
				}
			}
		}
	}
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
