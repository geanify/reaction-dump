package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"github.com/micmonay/keybd_event"
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

	str := "tree " + folder + " -f -i --filelimit 30 | grep -E '.jpg|png|gif|jpeg' | grep -E '" + text + "'"
	out := _execOutput(str)
	// fmt.Println(strings.Split(out, "\n"))
	return strings.Split(out, "\n")
}
func deferPaste(window fyne.Window) {
	// str := "sleep 1 ; xclip -selection clipboard -o >"
	window.Close()
	go func() {
		kb, err := keybd_event.NewKeyBonding()
		if err != nil {
			os.Exit(0)
		}
		kb.HasCTRL(true)
		kb.SetKeys(keybd_event.VK_V)
		// _exec(str)
		time.Sleep(500 * time.Millisecond)

		kb.Press()
		os.Exit(0)
	}()
}

func executeCall(th *Throttler) {
	for true {
		th.executeAfter()
	}
}
