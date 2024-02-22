package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
		log.Fatal(err)
	}
	outStr := string(out[:])
	return outStr
}

func copyImage(location string, filename string) {
	str := "cat " + location + filename + "| xclip -selection clipboard -target image/png -i"
	_exec(str)
}

func textLookUp(text string) {
	str := "tree -f -i | grep " + text
	out := _execOutput(str)
	fmt.Println(out)
}

func main() {
	// location := ""
	args := os.Args

	if len(os.Args[1:]) < 1 {
		copyImage("", "test.jpg")
	} else {
		copyImage(args[1], "test.jpg")
	}
	textLookUp("te")

	// exec.Command("touch test.txt")
}
