package tester

import (
	"net/http"
	"sync"
	"time"
	"tviewTest/settings"
)

type TestResult struct {
	IsUp   bool
	Url    string
	Name   string
	LastUp *time.Time
}

func Tester(yaml settings.YAML) (*sync.WaitGroup, []chan TestResult) {
	var wg sync.WaitGroup

	testResultChans := make([]chan TestResult, len(yaml.URLs))
	wg.Add(len(yaml.URLs))

	for i, Url := range yaml.URLs {
		url := Url
		testResCh := GoTest(testResultChans, url, yaml)
		testResultChans[i] = testResCh
	}
	return &wg, testResultChans
}

func GoTest(testResultChans []chan TestResult, Url settings.Url, yaml settings.YAML) chan TestResult {
	ch := make(chan TestResult)
	go func() {
		testResultChans = append(testResultChans, ch)
		for {
			_, err := http.Get(Url.Url)
			var testResult TestResult
			testResult.Url = Url.Url
			testResult.Name = Url.Name
			testResult.LastUp = nil

			if err != nil {
				testResult.IsUp = false
			} else {
				testResult.IsUp = true
			}

			// fmt.Println(url, "==>", testResult.IsUp)d

			ch <- testResult

			time.Sleep(time.Second * time.Duration(yaml.Interval_s))
		}
	}()

	return ch
}
