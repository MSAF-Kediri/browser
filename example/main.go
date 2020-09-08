package main

import (
	"browser"
	"fmt"
	"log"
	"os"
)

func main() {
	var br = new(browser.Browser)

	logName := "log_browser"
	f, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetFlags(0)
	log.SetOutput(f)

	br.LogFile = f
	br.Init(4444)
	defer br.QuitBrowser()

	// Resize Window
	_ = br.Wd.ResizeWindow("", 1600, 900)

	windowHandles, _ := br.Wd.WindowHandles()
	fmt.Println("Window Handles", windowHandles)

	// Navigate to specific url.
	url := "https://ibank.klikbca.com/"
	//url := "https://ibank.bni.co.id/"
	//url := "https://ibank.bankmandiri.co.id/"
	//url := "https://ib.bri.co.id/ib-bri/"
	if err := br.Wd.Get(url); err != nil {
		panic(err)
	}
}
