package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func createImage(path string, window fyne.Window) *fyne.Container {
	if len(path) < 1 {
		return nil
	}
	iconFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(iconFile)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	resource := fyne.NewStaticResource("icon", b)
	image := canvas.NewImageFromResource(resource)

	image.FillMode = 1

	btn := widget.NewButton("", func() {
		copyImage(path)
		window.Close()
		deferPaste(window)
	})
	btn.Alignment = 2

	content := container.New(layout.NewStackLayout(), btn, image)
	return content
}

func imageList(images []string, window fyne.Window) []*fyne.Container {
	imageList := make([]*fyne.Container, 0)
	for i := 0; i < len(images); i++ {
		newImage := createImage(images[i], window)
		if newImage != nil {
			imageList = append(imageList, newImage)
		}
	}
	return imageList
}

func handleImageLookUp(search string, window fyne.Window, content *fyne.Container) {
	results := textLookUp(search)
	if len(search) < 1 {
		return
	}
	imageList := imageList(results, window)

	for i := 0; i < len(imageList); i++ {
		if i < len(content.Objects) {
			content.Objects[i] = imageList[i]
		} else {
			content.Objects = append(content.Objects, imageList[i])
		}
	}

	for i := len(imageList); i < len(content.Objects); i++ {
		content.Objects[i] = nil
	}
	window.Canvas().Refresh(content)
}

func handleEnter(search string, window fyne.Window) {
	if len(search) < 1 {
		window.Close()
		return
	}
	results := textLookUp(search)

	if len(results) < 1 {
		return
	}
	path := results[0]

	copyImage(path)
	window.Close()
	deferPaste(window)

}

func render(window fyne.Window, th *Throttler) {
	search := widget.NewEntry()
	// text := canvas.NewText("Overlay", color.Black)
	// imgWidget := widget.NewCard("test", "test2", img)
	imageContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(190, 108)))
	handleImageLookUp("", window, imageContainer)
	content := container.New(layout.NewVBoxLayout(), search, imageContainer)
	search.PlaceHolder = "Your Face When..."
	search.OnChanged = func(s string) {
		th.reset()
		th.call = func() {
			handleImageLookUp(s, window, imageContainer)
		}
		th.shouldExecute = true
	}

	search.OnSubmitted = func(s string) {
		handleEnter(s, window)
	}

	window.Canvas().SetOnTypedKey(func(k *fyne.KeyEvent) {
		if k.Name == fyne.KeyEscape {
			window.Close()
		}
	})

	go executeCall(th)

	window.SetContent(content)
	window.Canvas().Focus(search)
}
