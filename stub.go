package main

import (
	"fmt"
	"github.com/sidav/golibrl/procedural_generation"
)

func main() {
	cave := procedural_generation.MakeCave(60, 20, 3, -1)
	for _, s := range *cave {
		fmt.Println(s)
	}
	//console.Init_console("test", console.TCellRenderer)
	//defer console.Close_console()
	//console.SetFgColor(console.WHITE)
	//for y, s := range *cave {
	//	console.PutString(s, 0, y)
	//}
	//console.Flush_console()
	//console.ReadKey()
}
