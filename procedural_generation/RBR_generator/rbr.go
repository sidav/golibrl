package RBR_generator

type RBR struct {
	tiles      [][]tile
	mapw, maph int
}

func (rbr *RBR) Init(w, h int) {
	rbr.tiles = make([][]tile, w)
	for row := range rbr.tiles {
		rbr.tiles[row] = make([]tile, h)
	}
	rbr.mapw = w
	rbr.maph = h
}

func (rbr *RBR) digSpace(x, y, w, h int) {
	for cx:=0; cx < x+w; cx++ {
		for cy:=0;cy<y+h; cy++ {
			rbr.tiles[cx][cy].tiletype = TFLOOR
		}
	}
}

func (rbr *RBR) Generate() {

}

func (rbr *RBR) GetMapChars() *[][]rune {
	runearr := make([][]rune, rbr.mapw)
	for row := range runearr {
		runearr[row] = make([]rune, rbr.maph)
	}
	for x:=0; x < rbr.mapw; x++ {
		for y:=0;y<rbr.maph; y++ {
			runearr[x][y] = rbr.tiles[x][y].toRune()
		}
	}
	return &runearr
}
