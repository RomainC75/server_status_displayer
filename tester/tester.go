package tester

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"tviewTest/settings"
)

type TestingResult struct {
	isUp bool
	url  string
}

func Tester(yaml settings.YAML) *sync.WaitGroup {
	var wg sync.WaitGroup

	testResultChans := make([]chan TestingResult, len(yaml.URLs))
	wg.Add(len(yaml.URLs))

	for i, url := range yaml.URLs {
		url := url
		testResCh := GoTest(testResultChans, url, yaml)
		testResultChans[i] = testResCh
	}
	return &wg
}

func GoTest(testResultChans []chan TestingResult, url string, yaml settings.YAML) chan TestingResult {
	ch := make(chan TestingResult)
	go func() {
		// var ch chan TestingResult
		testResultChans = append(testResultChans, ch)
		for {
			_, err := http.Get(url)
			var testingResult TestingResult
			testingResult.url = url

			if err != nil {
				testingResult.isUp = false
			} else {
				testingResult.isUp = true
			}

			fmt.Println(url, "==>", testingResult.isUp)

			ch <- testingResult

			time.Sleep(time.Second * time.Duration(yaml.Interval_s))
		}
	}()

	return ch
}
