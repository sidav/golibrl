package BSP_generator

import "github.com/sidav/golibrl/random/additive_random"

type container struct {
	x, y, w, h int
}

func (r *container) returnCenter() (int, int) {
	return (r.x + r.w/2), (r.y + r.h/2)
}

type treeNode struct {
	parent, left, right *treeNode
	room                *container
}

type ReturningMap struct { //this struct is returned from generation routine.
	dmap []rune
}

func (m *ReturningMap) init() {
	m.dmap = make([]rune, MAP_W*MAP_H)
	for i := 0; i < len(m.dmap); i++ {
		m.dmap[i] = FLOOR
	}
}

func (m *ReturningMap) GetCell(x, y int) rune {
	return m.dmap[x+MAP_W*y]
}

func (m *ReturningMap) SetCell(cell rune, x, y int) {
	m.dmap[x+MAP_W*y] = cell
}

func (m *ReturningMap) CountWallsAround(x, y int) int {
	sum := 0
	if m.dmap[(x-1)+MAP_W*y] != FLOOR {
		sum++
	}
	if m.dmap[(x+1)+MAP_W*y] != FLOOR {
		sum++
	}
	if m.dmap[x+MAP_W*(y-1)] != FLOOR {
		sum++
	}
	if m.dmap[x+MAP_W*(y+1)] != FLOOR {
		sum++
	}
	return sum
}

func (m *ReturningMap) CountDoorsAround(x, y int) int {
	sum := 0
	if m.dmap[(x-1)+MAP_W*y] == DOOR {
		sum++
	}
	if m.dmap[(x+1)+MAP_W*y] == DOOR {
		sum++
	}
	if m.dmap[x+MAP_W*(y-1)] == DOOR {
		sum++
	}
	if m.dmap[x+MAP_W*(y+1)] == DOOR {
		sum++
	}
	return sum
}

func getSplitRangeForPercent(wh int, percent int) (int, int) {
	min := wh * percent / 100
	return min, wh - min
}

func (t *treeNode) splitHoriz() { // splits node into "lower" and "upper"
	current_w := t.room.w
	current_h := t.room.h
	current_x := t.room.x
	current_y := t.room.y
	minSplSize, maxSplSize := getSplitRangeForPercent(current_h, SPLIT_MIN_RATIO)
	// Let's try to split the node without breaking min room size constraints
	for try := 0; try < TRIES_FOR_SPLITTING; try++ {
		upper_h := rnd.RandInRange(minSplSize, maxSplSize)
		lower_h := current_h - upper_h + 1
		if upper_h < MIN_ROOM_H || lower_h < MIN_ROOM_H {
			continue
		} else { // Okay, sizes are acceptable. Let's do the split
			upperNode := treeNode{parent: t, room: &container{x: current_x, y: current_y, w: current_w, h: upper_h}}
			// Most error-probable place:
			lowerNode := treeNode{parent: t, room: &container{x: current_x, y: current_y + upper_h - 1, w: current_w, h: lower_h}}
			// hm... Left is upper and right is lower. Everything is obvious.
			t.left = &upperNode
			t.right = &lowerNode
			return
		}
	}
}

func (t *treeNode) splitVert() { // splits node into left and right
	current_w := t.room.w
	current_h := t.room.h
	current_x := t.room.x
	current_y := t.room.y
	minSplSize, maxSplSize := getSplitRangeForPercent(current_w, SPLIT_MIN_RATIO)
	// Let's try to split the node without breaking min room size constraints
	for try := 0; try < TRIES_FOR_SPLITTING; try++ {
		left_w := rnd.RandInRange(minSplSize, maxSplSize)
		right_w := current_w - left_w + 1
		if left_w < MIN_ROOM_W || right_w < MIN_ROOM_W {
			continue
		} else { // Okay, sizes are acceptable. Let's do the split
			leftNode := treeNode{parent: t, room: &container{x: current_x, y: current_y, w: left_w, h: current_h}}
			// Most error-probable place:
			rightNode := treeNode{parent: t, room: &container{x: current_x + left_w - 1, y: current_y, w: right_w, h: current_h}}
			t.left = &leftNode
			t.right = &rightNode
			return
		}
	}
}

func (t *treeNode) splitNTimes(n int) {
	if n == 0 {
		return
	}
	toSplitOrNotToSplit := rnd.Rand(100)
	if toSplitOrNotToSplit < SPLIT_PROBABILITY || t.room.w > MAX_ROOM_W || t.room.h > MAX_ROOM_H {
		horOrVert := rnd.Rand(100)
		if horOrVert < HORIZ_PROBABILITY {
			t.splitHoriz()
		} else {
			t.splitVert()
		}
		if t.left != nil && t.right != nil { //if split was successful
			t.left.splitNTimes(n - 1)
			t.right.splitNTimes(n - 1)
		}
	}
}

func countOutsizedRooms(node *treeNode) int {
	total := 0
	if node.left == nil {
		if node.room.w > MAX_ROOM_W || node.room.h > MAX_ROOM_H {
			return 1
		}
		return 0
	} else {
		total += countOutsizedRooms(node.left)
		total += countOutsizedRooms(node.right)
	}
	return total
}

/////////////////////////////////////////

const (
	WALL                 = '#'
	RIVER                = '~'
	DOOR                 = '+'
	FLOOR                = '.'
	TRIES_FOR_SPLITTING  = 10
	TRIES_FOR_GENERATION = 1000
	MAX_OUTSIZED_ROOMS   = 5
)

var (
	MAP_W, MAP_H      int
	treeRoot          *treeNode
	SPLIT_PROBABILITY = 70 // in percent.
	SPLIT_MIN_RATIO   = 30 // in percent.
	MIN_ROOM_W        = 5
	MIN_ROOM_H        = 5
	MAX_ROOM_W        = 20 // this and next lines are not guaranteed. Think of them as a recommendations.
	MAX_ROOM_H        = 10 //
	HORIZ_PROBABILITY = 30 // in percent. Horiz splits should occur less frequently than vertical ones because of w > h
	rnd = additive_random.FibRandom{}
)

func GenerateDungeon(width, height, splits, sp_prob, sp_ratio, h_prob, riverWidth int) *ReturningMap {
	rnd.InitDefault()
	MAP_W = width
	MAP_H = height
	if splits == 0 {
		splits = 6
	}
	if sp_prob != 0 {
		SPLIT_PROBABILITY = sp_prob
	}
	if sp_ratio != 0 {
		SPLIT_MIN_RATIO = sp_ratio
	}
	if h_prob != 0 {
		HORIZ_PROBABILITY = h_prob
	}

	for i := 0; i < TRIES_FOR_GENERATION; i++ {
		// generate parent node
		treeRoot = &treeNode{room: &container{x: 0, y: 0, w: MAP_W, h: MAP_H}}
		// recursively split into rooms
		treeRoot.splitNTimes(splits)
		if countOutsizedRooms(treeRoot) > MAX_OUTSIZED_ROOMS {
			continue
		} else {
			break
		}
	}

	// init returning struct
	result := &ReturningMap{}
	result.init()

	renderTreeToDungeonMap(treeRoot, result)
	if riverWidth > 0 {
		addRiverForDungeonMap(result, riverWidth)
	}
	addDoorsForDungeonMap(treeRoot, result)

	return result
}

func renderTreeToDungeonMap(node *treeNode, dmap *ReturningMap) {
	// recursively traverse through nodes and draw their containers
	if node.left != nil {
		renderTreeToDungeonMap(node.left, dmap)
		renderTreeToDungeonMap(node.right, dmap)
		return
	}
	for x := node.room.x; x < node.room.x+node.room.w; x++ {
		dmap.SetCell(WALL, x, node.room.y)
		dmap.SetCell(WALL, x, node.room.y+node.room.h-1)
	}
	for y := node.room.y; y < node.room.y+node.room.h; y++ {
		dmap.SetCell(WALL, node.room.x, y)
		dmap.SetCell(WALL, node.room.x+node.room.w-1, y)
	}
}

func addRiverForDungeonMap(dmap *ReturningMap, riverWidth int) {
	x := rnd.RandInRange(MAP_W/3, MAP_W*2/3)
	bridgeHeight := 2
	bridgeYCoord := rnd.RandInRange(1, MAP_H-1-bridgeHeight)
	for y := 0; y < MAP_H; y++ {
		dmap.SetCell(FLOOR, x-1, y)
		dmap.SetCell(FLOOR, x-2, y)
		dmap.SetCell(FLOOR, x+riverWidth, y)
		dmap.SetCell(FLOOR, x+riverWidth+1, y)
		for cx := 0; cx < riverWidth; cx++ {
			if y >= bridgeYCoord && y < bridgeYCoord+bridgeHeight {
				dmap.SetCell(FLOOR, x+cx, y)
			} else {
				dmap.SetCell(RIVER, x+cx, y)
			}
		}
		leftOrRight := rnd.RandInRange(0, 2)
		if leftOrRight == 0 {
			x--
		}
		if leftOrRight == 1 {
			x++
		}
	}
}

// BUGGED! Rooms connectivity still not guaranteed!
func addDoorsForDungeonMap(node *treeNode, dmap *ReturningMap) {
	if node.left != nil {
		lx, ly := node.left.room.returnCenter()
		rx, ry := node.right.room.returnCenter()

		if ly == ry {
			// ly += randInRange(-MIN_ROOM_H/2, MIN_ROOM_H/2)
			for x := lx; x < rx; x++ {
				if dmap.GetCell(x, ly) == WALL {
					if dmap.CountWallsAround(x, ly) > 2 {
						continue
					}
					dmap.SetCell(DOOR, x, ly)
				}
			}
		}
		if lx == rx {
			// lx += randInRange(-MIN_ROOM_W/2, MIN_ROOM_W/2)
			for y := ly; y < ry; y++ {
				if dmap.GetCell(lx, y) == WALL {
					if dmap.CountWallsAround(lx, y) > 2 {
						continue
					}
					dmap.SetCell(DOOR, lx, y)
				}
			}
		}
		addDoorsForDungeonMap(node.left, dmap)
		addDoorsForDungeonMap(node.right, dmap)
	}
}
