package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Print("error", err)
	}
	defer sdl.Quit()

	fmt.Print("Dengueboy")
}
