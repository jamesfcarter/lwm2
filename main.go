package main

import (
	"log"

	"github.com/jamesfcarter/lwm2/wm"
)

func main() {
	wm, err := wm.Init()
	if err != nil {
		log.Fatalf("failed to init: %v", err)
	}
	wm.Run()
}
