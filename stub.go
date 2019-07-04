package main

import (
	"github.com/sidav/golibrl/console"
	"github.com/sidav/golibrl/procedural_generation/Fractal_landscape"
)

func main() {
	//cave := CA_cave.MakeCave(60, 20, 3, -1)
	//for _, s := range *cave {
	//	fmt.Println(s)
	//}

	land := Fractal_landscape.GenHeightMap(129, 65)
	// return

	console.Init_console("test", console.TCellRenderer)
	defer console.Close_console()
	console.SetFgColor(console.WHITE)
	for i := 0; i < len(*land); i++ {
		str := ' '
		for j := 0; j < len((*land)[0]); j++ {
			switch cur := (*land)[i][j]; {
			case cur < -10:
				str = '~'
				console.SetFgColor(console.DARK_BLUE)
			case cur < 0:
				str = '~'
				console.SetFgColor(console.BLUE)
			case cur < 9:
				str = '.'
				console.SetFgColor(console.YELLOW)
			case cur < 22:
				str = ','
				console.SetFgColor(console.DARK_YELLOW)
			case cur < 40:
				str = 'T'
				console.SetFgColor(console.GREEN)
			case cur < 50:
				str = '^'
				console.SetFgColor(console.DARK_GRAY)
			default:
				str = '^'
				console.SetFgColor(console.WHITE)
			}
			console.PutChar(str, i, j)
		}
	}
	console.Flush_console()
	console.ReadKey()
}
