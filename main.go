package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Throttler struct {
	start         time.Time
	shouldExecute bool
	call          func()
}

func (th *Throttler) executeAfter() {
	if !th.shouldExecute {
		return
	}
	now := time.Now()
	elapsed := now.Sub(th.start).Milliseconds()
	if elapsed < 300 {
		time.Sleep(300 * time.Millisecond)
	}
	th.shouldExecute = false
	th.call()
	th.reset()
}

func (th *Throttler) reset() {
	th.start = time.Now()
}

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
	str := "cat '" + path + "' | xclip -selection clipboard -target image/png -i"
	_exec(str)
}

func textLookUp(text string) []string {
	args := os.Args
	folder := "."

	if len(os.Args[1:]) >= 1 {
		folder = args[1]
	}

	str := "tree " + folder + " -f -i | grep -E '.jpg|png|gif|jpeg' | grep -E '" + text + "'"
	out := _execOutput(str)
	// fmt.Println(strings.Split(out, "\n"))
	return strings.Split(out, "\n")
}

func createImage(path string) *fyne.Container {
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
		os.Exit(0)
	})
	btn.Alignment = 2

	content := container.New(layout.NewStackLayout(), btn, image)
	return content
}

func imageList(images []string) []*fyne.Container {
	imageList := make([]*fyne.Container, 0)
	for i := 0; i < len(images); i++ {
		newImage := createImage(images[i])
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
	imageList := imageList(results)

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
	handleRefresh(window, content)
}
func handleUpdate(search string, window fyne.Window, content *fyne.Container) {
	handleImageLookUp(search, window, content)
}

func handleRefresh(window fyne.Window, content *fyne.Container) {
	window.Canvas().Refresh(content)
}

func render(window fyne.Window, th *Throttler) {
	search := widget.NewEntry()
	// text := canvas.NewText("Overlay", color.Black)
	// imgWidget := widget.NewCard("test", "test2", img)
	imageContainer := container.New(layout.NewGridWrapLayout(fyne.NewSize(190, 108)))
	handleUpdate("", window, imageContainer)
	content := container.New(layout.NewVBoxLayout(), search, imageContainer)

	search.OnChanged = func(s string) {
		th.reset()
		th.call = func() {
			handleUpdate(s, window, imageContainer)
		}
		th.shouldExecute = true
	}

	go executeCall(th)

	// mx, err := fyne.Vector2{10000.0, 50.0}
	// search.Size().Max(mx)
	window.SetContent(content)
	window.Canvas().Focus(search)
}

func executeCall(th *Throttler) {
	for true {
		th.executeAfter()
	}
}

func main() {
	th := &Throttler{start: time.Now(), call: func() {}}
	myApp := app.New()
	window := myApp.NewWindow("Reaction Dump")
	render(window, th)
	go executeCall(th)

	window.CenterOnScreen()
	window.SetFixedSize(true)
	window.Resize(fyne.NewSize(800, 600))
	window.ShowAndRun()

	// exec.Command("touch test.txt")
}
