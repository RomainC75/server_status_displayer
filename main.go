package main

import (
	"log"
	"time"
	"tviewTest/display"
	"tviewTest/settings"
	"tviewTest/tester"
)

type Data struct {
	Date   time.Time
	Region string
	Price  float32
}

func main() {
	err := settings.SetSettings()
	if err != nil {
		log.Fatal(err.Error())
	}

	yaml := settings.Settings.Get()

	wg, testResultChans := tester.Tester(yaml)

	display := display.NewDisplayer(yaml)
	display.GoMerger(testResultChans)

	display.Draw()

	wg.Wait()
}
