package main

import (
	"github.com/nsf/termbox-go" // Consider termui instead
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic("Could not initialize termbox.")
	}

	termbox.Clear(termbox.ColorYellow, termbox.ColorYellow)
	termbox.SetCell(2, 2, 'H', termbox.ColorBlue, termbox.ColorBlack)
	termbox.PollEvent()

	defer termbox.Close()
}
