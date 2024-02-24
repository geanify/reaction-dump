package main

import (
	"os"
	"time"

	"github.com/micmonay/keybd_event"
)

func autoPaste() {
	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		os.Exit(0)
	}
	kb.HasCTRL(true)
	kb.SetKeys(keybd_event.VK_V)
	// _exec(str)
	time.Sleep(500 * time.Millisecond)

	kb.Press()
}
