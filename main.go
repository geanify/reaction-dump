package main

import (
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

func main() {
	th := &Throttler{start: time.Now(), call: func() {}}
	myApp := app.New()
	window := myApp.NewWindow("Reaction Dump")
	render(window, th)
	go executeCall(th)

	window.SetOnClosed(func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	})
	window.CenterOnScreen()
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()

	for true {
		time.Sleep(300 * time.Second)
	}

	// exec.Command("touch test.txt")
}
