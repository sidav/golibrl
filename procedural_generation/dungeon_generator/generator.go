package dungeon_generator

import (
	rnd "github.com/sidav/golibrl/random"
	"github.com/sidav/golibrl/random/additive_random"
	"sort"
)

var random = additive_random.FibRandom{}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func arrIntersection(a, b []int) bool {
	for ai := range a {
		for bi := range b {
			if a[ai] == b[bi] {
				return true
			}
		}
	}
	return false
}

type Generator struct {
	Width, Height, MaxRooms, MinRoomXY,
	MaxRoomXY, RandomConnections, RandomSpurs int
	RoomsOverlap bool

	level [][]rune
	corridor_list [][][]int
	room_list []*[]int
}

func (g *Generator) Init() {
	random.InitDefault()
	g.corridor_list = make([][][]int, 0)
}


func (g *Generator) appendCorridors(corridors *[][]int) {
	g.corridor_list = append(g.corridor_list, (*corridors))
}

func (g *Generator) appendRoom(room *[]int) {
	g.room_list = append(g.room_list, room)
}

func (g *Generator) gen_room() []int {
	w := random.RandInRange(g.MinRoomXY, g.MaxRoomXY)
	h := random.RandInRange(g.MinRoomXY, g.MaxRoomXY)
	x := random.RandInRange(1, (g.Width - w - 1))
	y := random.RandInRange(1, (g.Height - h - 1))

	return []int {x, y, w, h}
}

func (g *Generator) roomsOverlapping(room *[]int, room_list *[]*[]int) bool {
	x := (*room)[0]
	y := (*room)[1]
	w := (*room)[2]
	h := (*room)[3]
	for _, current_room := range *room_list {
		// The rectangles don't overlap if
		// one rectangle's minimum in some dimension
		// is greater than the other's maximum in
		// that dimension.
		if x < ((*current_room)[0] + (*current_room)[2]) &&
			(*current_room)[0] < (x + w) && 	
			y < ((*current_room)[1] + (*current_room)[3]) &&	
			(*current_room)[1] < (y + h) {
			return true 
		}
	}
	return false
}

func (g *Generator) corridorBetweenPoints(x1, y1, x2, y2 int, joinType string) [][]int {
	if joinType == "" {
		joinType = "either"
	}
	if x1 == x2 || y1 == y2 {
		return [][]int {{x1, y1}, {x2, y2}}
	} else {
		join := "none"
		if joinType == "either" && arrIntersection([]int{0, 1}, []int{x1, y1, x2, y2}) {
			join = "bottom"
		} else if joinType == "either" && arrIntersection([]int{g.Width -1, g.Width -2}, []int{x1, x2}) ||
			arrIntersection([]int{g.Height -1, g.Height -2}, []int{y1, y2}) {
			join = "top"
		} else if joinType == "either" {
			join = []string{"top", "bottom"}[rnd.Random(2)]
		} else {
			join = joinType
		}

		if join == "top" {
			return [][]int{{x1, y1}, {x1, y2}, {x2, y2}}
		} else if join == "bottom" {
			return [][]int{{x1, y1}, {x2, y1}, {x2, y2}}
		}
		panic("No nil case?")
	}
}

func (g *Generator) joinRooms(room1, room2 []int, joinType string) {
	if joinType == "" {
		joinType = "either"
	}
	sortedRoom := [][]int{room1, room2}
	// WARNING: maybe change > to < ? 
	if room1[0] > room2[0] {
		sortedRoom = [][]int{room2, room1}
	}
	x1 := sortedRoom[0][0]
	y1 := sortedRoom[0][1]
	w1 := sortedRoom[0][2]
	h1 := sortedRoom[0][3]
	x1_2 := x1 + w1 - 1
	y1_2 := y1 + h1 - 1

	x2 := sortedRoom[1][0]
	y2 := sortedRoom[1][1]
	w2 := sortedRoom[1][2]
	h2 := sortedRoom[1][3]
	x2_2 := x2 + w2 - 1
	y2_2 := y2 + h2 - 1

	if x1 < (x2 + w2) && x2 < (x1 + w1) {
		jx1 := random.RandInRange(x2, x1_2)
		jx2 := jx1
		tmp_y := []int{y1, y2, y1_2, y2_2}
		// maybe reverse sort?
		sort.Ints(tmp_y)
		jy1 := tmp_y[1] + 1
		jy2 := tmp_y[2] - 1

		corridors := g.corridorBetweenPoints(jx1, jy1, jx2, jy2, "")
		g.appendCorridors(&corridors)

	} else if y1 < (y2 + h2) && y2 < (y1 + h1) {
		var jy1, jy2 int
		if y2 > y1 {
			jy1 = random.RandInRange(y2, y1_2)
			jy2 = jy1
		} else {
			jy1 = random.RandInRange(y1, y2_2)
			jy2 = jy1
		}
		tmp_x := []int{x1, x2, x1_2, x2_2}
		sort.Ints(tmp_x)
		jx1 := tmp_x[1] + 1
		jx2 := tmp_x[2] - 1

		corridors := g.corridorBetweenPoints(jx1, jy1, jx2, jy2, "")
		g.appendCorridors(&corridors)
	} else {
		join := "None"
		if joinType == "either" {
			join = []string{"top", "bottom"}[rnd.Random(2)]
		} else {
			join = joinType
		}

		if join == "top" {
			if y2 > y1 {
				jx1 := x1_2 + 1
				jy1 := random.RandInRange(y1, y1_2)
				jx2 := random.RandInRange(x2, x2_2)
				jy2 := y2 - 1
				corridors := g.corridorBetweenPoints(
					jx1, jy1, jx2, jy2, "bottom")
				g.appendCorridors(&corridors)
			} else {
				jx1 := random.RandInRange(x1, x1_2)
				jy1 := y1 - 1
				jx2 := x2 - 1
				jy2 := random.RandInRange(y2, y2_2)
				corridors := g.corridorBetweenPoints(
					jx1, jy1, jx2, jy2, "top")
				g.appendCorridors(&corridors)
			}

		} else if join == "bottom" {
			if y2 > y1 {
				jx1 := random.RandInRange(x1, x1_2)
				jy1 := y1_2 + 1
				jx2 := x2 - 1
				jy2 := random.RandInRange(y2, y2_2)
				corridors := g.corridorBetweenPoints(
					jx1, jy1, jx2, jy2, "top")
				g.appendCorridors(&corridors)
			} else {
				jx1 := x1_2 + 1
				jy1 := random.RandInRange(y1, y1_2)
				jx2 := random.RandInRange(x2, x2_2)
				jy2 := y2_2 + 1
				corridors := g.corridorBetweenPoints(
					jx1, jy1, jx2, jy2, "bottom")
				g.appendCorridors(&corridors)
			}
		}
	}
}

func (self *Generator) GenLevel() [][]rune {
	self.level = make([][]rune, self.Height)
	for i := 0; i < self.Height; i++ {
		self.level[i] = make([]rune, self.Width)
	}
	max_iters := self.MaxRooms * 5
	for a := 0; a < max_iters; a++ {
		tmp_room := self.gen_room()

		if self.RoomsOverlap || len(self.room_list) == 0 {
			self.appendRoom(&tmp_room)
		} else {
			tmp_room = self.gen_room()
			if !self.roomsOverlapping(&tmp_room, &self.room_list) {
				self.appendRoom(&tmp_room)
			}
		}
		if len(self.room_list) >= self.MaxRooms {
			break
		}
	}
	// connect rooms
	for a := 0; a < len(self.room_list) - 1; a++ {
		self.joinRooms(*self.room_list[a], *self.room_list[a+1], "")
	}
	// do the random joins
	for a := 0; a < self.RandomConnections; a++ {
		self.joinRooms(
			*self.room_list[random.Rand(len(self.room_list))],
			*self.room_list[random.Rand(len(self.room_list))],
			"")
	}
	// do the spurs
	for a := 0; a < self.RandomSpurs; a++ {
		room_1 := []int{random.RandInRange(2, self.Width- 2), random.RandInRange(
			2, self.Height- 2), 1, 1}
		room_2 := self.room_list[random.Rand(len(self.room_list))]
		self.joinRooms(room_1, *room_2,"")
	}
	// fill the map
	for room_num := range self.room_list {
		room := *self.room_list[room_num]
		for b := 0; b < room[2]; b++ {
			for c := 0; c < room[3]; c++ {
				self.level[room[1] + c][room[0] + b] = '.'
			}
		}
	}
	for _, corridor := range self.corridor_list {
		x1 := corridor[0][0]
		y1 := corridor[0][1]
		x2 := corridor[1][0]
		y2 := corridor[1][1]
		for w := 0; w < (abs(x1 - x2) + 1); w++ {
			for h := 0; h < (abs(y1 - y2) + 1); h++ {
				self.level[min(y1, y2) + h][min(x1, x2) + w] = '.'
			}
		}
		if len(corridor) == 3 {
			x3 := corridor[2][0]
			y3 := corridor[2][1]
			for w := 0; w < (abs(x2-x3) + 1); w++ {
				for h := 0; h < (abs(y2-y3) + 1); h++ {
					self.level[min(y2, y3)+h][min(x2, x3)+w] = '.'
				}
			}
		}
	}
	return self.level
}
