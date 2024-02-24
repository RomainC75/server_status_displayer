package tester

import (
	"fmt"
	"net/http"
	"sync"
	"time"
	"tviewTest/settings"
)

type TestResult struct {
	IsUp bool
	Url  string
}

func Tester(yaml settings.YAML) (*sync.WaitGroup, []chan TestResult) {
	var wg sync.WaitGroup

	testResultChans := make([]chan TestResult, len(yaml.URLs))
	wg.Add(len(yaml.URLs))

	for i, url := range yaml.URLs {
		url := url
		testResCh := GoTest(testResultChans, url, yaml)
		testResultChans[i] = testResCh
	}
	return &wg, testResultChans
}

func GoTest(testResultChans []chan TestResult, url string, yaml settings.YAML) chan TestResult {
	ch := make(chan TestResult)
	go func() {
		testResultChans = append(testResultChans, ch)
		for {
			_, err := http.Get(url)
			var testResult TestResult
			testResult.Url = url

			if err != nil {
				testResult.IsUp = false
			} else {
				testResult.IsUp = true
			}

			fmt.Println(url, "==>", testResult.IsUp)

			ch <- testResult

			time.Sleep(time.Second * time.Duration(yaml.Interval_s))
		}
	}()

	return ch
}
