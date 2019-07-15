package epfov

type line_t struct {
	xi, yi, xf, yf int
}

type viewbump_t struct {
	x,y int
	refcount int
	parent *viewbump_t
}

type view_t struct {
	shallow_line line_t
	steep_line line_t
	shallow_bump, steep_bum *viewbump_t
}


