package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func _exec(str string) {
	cmd := exec.Command("bash", "-c", str)
	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func _execOutput(str string) string {
	cmd := exec.Command("bash", "-c", str)
	out, err := cmd.Output()

	if err != nil {
		return ""
	}
	outStr := string(out[:])
	return outStr
}

func copyImage(path string) {
	str := "cat " + path + "| xclip -selection clipboard -target image/png -i"
	_exec(str)
}

func textLookUp(text string) {
	str := "tree -f -i | grep " + text
	out := _execOutput(str)
	fmt.Println(out)
}

func createImage(path string) *widget.Button {
	iconFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(iconFile)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	btn := widget.NewButtonWithIcon("Browse", fyne.NewStaticResource("icon", b), func() {
		// do something
		copyImage(path)
	})
	return btn
}

func imageList(images []string) []*widget.Button {
	imageList := make([]*widget.Button, 0)
	for i := 0; i < len(images); i++ {
		newImage := createImage("reactions/test.jpg")
		imageList = append(imageList, newImage)
	}
	return imageList
}

func render(window fyne.Window) {
	search := widget.NewEntry()
	search.OnChanged = textLookUp

	img := createImage("reactions/test.jpg")
	img2 := createImage("reactions/test.jpg")
	// text := canvas.NewText("Overlay", color.Black)
	// imgWidget := widget.NewCard("test", "test2", img)

	content := container.New(layout.NewGridLayoutWithRows(3), search, img2, img, img)

	window.SetContent(content)
}

func main() {
	args := os.Args

	if len(os.Args[1:]) < 1 {
		copyImage("test.jpg")
	} else {
		copyImage(args[1] + "test.jpg")
	}
	textLookUp("te")

	myApp := app.New()
	window := myApp.NewWindow("Reaction Dump")
	render(window)

	window.CenterOnScreen()
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()

	// exec.Command("touch test.txt")
}
