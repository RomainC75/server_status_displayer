package main

import (
	"log"
	"time"
	"tviewTest/display"
	"tviewTest/settings"
	"tviewTest/tester"

	"github.com/k0kubun/pp"
)

type Data struct {
	Date   time.Time
	Region string
	Price  float32
}

func main() {
	// data := []Data{
	// 	{time.Now(), "US", 32.2},
	// 	{time.Now(), "France", 32.3},
	// 	{time.Now(), "Costarica", 32.1},
	// }

	err := settings.SetSettings()
	if err != nil {
		log.Fatal(err.Error())
	}

	yaml := settings.Settings.Get()

	pp.Print(yaml)

	wg, testResultChans := tester.Tester(yaml)

	display := display.NewDisplayer(yaml)
	display.GoMerger(testResultChans)

	display.Draw()

	wg.Wait()
}
