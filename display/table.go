package display

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"tviewTest/file"
	"tviewTest/settings"
	"tviewTest/tester"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Displayer struct {
	testResults    []tester.TestResult
	historyResults []file.History
	Interval_s     int
	app            *tview.Application
	flex           *tview.Flex
	stateTable     *tview.Table
	historyTable   *tview.Table
	m              sync.Mutex
}

func NewDisplayer(yaml settings.YAML) *Displayer {
	tr := []tester.TestResult{}
	for _, v := range yaml.URLs {
		tr = append(tr, tester.TestResult{
			Url:  v.Url,
			Name: v.Name,
		})
	}

	return &Displayer{
		testResults:    tr,
		historyResults: []file.History{},
		Interval_s:     yaml.Interval_s,
		app:            tview.NewApplication(),
		flex:           tview.NewFlex().SetDirection(tview.FlexColumn),
		stateTable:     tview.NewTable(),
		historyTable:   tview.NewTable(),
		m:              sync.Mutex{},
	}
}

func (d *Displayer) Draw() {
	go func() {
		i := 0
		for {
			d.app.QueueUpdateDraw(func() {
				d.RefreshStateTable()
				d.RefreshHistory()
				i++
			})
			time.Sleep(1 * time.Second)
		}

	}()
	d.flex.AddItem(d.stateTable, 0, 1, false).AddItem(d.historyTable, 0, 1, false)
	d.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlQ {
			fmt.Println("lkjsdf")
			file.CreateFile(d.historyResults)
			time.Sleep(time.Second * 4)
			// d.app.Stop()
			// => terminate every goRoutine
		}
		return event
	})
	if err := d.app.SetRoot(d.flex, true).Run(); err != nil {
		panic(err)
	}
}

func (d *Displayer) RefreshHistory() {
	d.historyTable.Clear()
	for row := 0; row < len(d.historyResults); row++ {

		// name

		d.historyTable.SetCell(row, 0, tview.NewTableCell(d.historyResults[row].Name).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignLeft))

		// start
		d.historyTable.SetCell(row, 1, tview.NewTableCell("[yellow]"+d.historyResults[row].Start.Format("2006-01-02 15:04:05")).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignLeft))

		// elapsed time
		elapsed_ns := file.GetElapsedTime_ns(d.historyResults[row])

		d.historyTable.SetCell(row, 2, tview.NewTableCell(fmt.Sprintf("[blue]%s", elapsed_ns)).
			SetTextColor(tview.Styles.PrimaryTextColor).
			SetAlign(tview.AlignRight))
	}
}

func (d *Displayer) RefreshStateTable() {
	d.stateTable.Clear()
	for row := 0; row < len(d.testResults); row++ {
		for column := 0; column < 3; column++ {
			align := tview.AlignRight

			var txt string
			switch column {
			case 0:
				txt = d.testResults[row].Name
				align = tview.AlignLeft
			case 1:
				if d.testResults[row].IsUp {
					txt = "[green]O"
				} else {
					txt = "[red]X"
				}
			case 2:
				if d.testResults[row].LastUp != nil {
					txt = "[yellow]" + d.testResults[row].LastUp.Format("2006-01-02 15:04:05")
				} else {
					txt = ""
				}
			}

			d.stateTable.SetCell(row, column,
				tview.NewTableCell(txt).
					SetTextColor(tview.Styles.PrimaryTextColor).
					SetAlign(align))
		}
	}
}

func (d *Displayer) GoMerger(channels []chan tester.TestResult) {
	// TODO : WAIT GROUP
	for i, ch := range channels {
		index := i
		channel := ch
		go func() {
			for {
				// select {
				// done
				// case result := <-channel:
				result := <-channel

				isNowUp := rand.Intn(3) < 1 && result.IsUp

				if d.testResults[index].IsUp && !isNowUp {
					now := time.Now()
					d.m.Lock()
					d.testResults[index].LastUp = &now
					d.m.Unlock()
				} else if isNowUp {
					if d.testResults[index].LastUp != nil {
						d.createHistoryEntry(d.testResults[index])
					}
					d.m.Lock()
					d.testResults[index].LastUp = nil
					d.m.Unlock()
				}
				d.m.Lock()
				d.testResults[index].IsUp = isNowUp
				d.m.Unlock()
			}
		}()
	}

}

func (d *Displayer) createHistoryEntry(tr tester.TestResult) {
	start := *tr.LastUp
	end := time.Now()
	history := file.History{
		Name:  tr.Name,
		Url:   tr.Url,
		Start: &start,
		End:   &end,
	}
	d.m.Lock()
	d.historyResults = append(d.historyResults, history)
	d.m.Unlock()
}
