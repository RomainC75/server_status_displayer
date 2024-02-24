package display

import (
	"math/rand"
	"sync"
	"time"
	"tviewTest/settings"
	"tviewTest/tester"

	"github.com/rivo/tview"
)

type Displayer struct {
	testResults []tester.TestResult
	Interval_s  int
	app         *tview.Application
	table       *tview.Table
	m           sync.Mutex
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
		testResults: tr,
		Interval_s:  yaml.Interval_s,
		app:         tview.NewApplication(),
		table:       tview.NewTable(),
		m:           sync.Mutex{},
	}
}

func (d *Displayer) Draw() {

	go func() {
		i := 0
		for {
			d.app.QueueUpdateDraw(func() {
				d.table.Clear()
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

						d.table.SetCell(row, column,
							tview.NewTableCell(txt).
								SetTextColor(tview.Styles.PrimaryTextColor).
								SetAlign(align))
					}
				}
				i++
			})
			time.Sleep(1 * time.Second)
		}
	}()

	if err := d.app.SetRoot(d.table, true).Run(); err != nil {
		panic(err)
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

				isUp := rand.Intn(3) < 1 && result.IsUp
				d.m.Lock()

				if d.testResults[index].IsUp && !isUp {
					now := time.Now()
					d.testResults[index].LastUp = &now
				}
				if isUp {
					d.testResults[index].LastUp = nil
				}
				d.testResults[index].IsUp = isUp
				d.m.Unlock()

				// }

			}
		}()
	}
}

// func display(data []Data) *tview.Table {
// 	table := tview.NewTable().SetFixed(1, 1)
// 	for row := 0; row < len(data); row++ {
// 		for column := 0; column < 3; column++ {
// 			color := tcell.ColorWhite
// 			if row == 0 {
// 				color = tcell.ColorYellow
// 			} else if column == 0 {
// 				color = tcell.ColorDarkCyan
// 			}
// 			align := tview.AlignLeft
// 			if row == 0 {
// 				align = tview.AlignCenter
// 			} else if column == 0 || column >= 4 {
// 				align = tview.AlignRight
// 			}

// 			var txt string
// 			switch column {
// 			case 0:
// 				txt = data[row].Date.String()
// 			case 1:
// 				txt = data[row].Region
// 			case 2:
// 				txt = fmt.Sprintf("%f", data[row].Price)
// 			}
// 			table.SetCell(row, column, &tview.TableCell{
// 				Text:  txt,
// 				Color: color,
// 				Align: align,
// 			})
// 		}
// 	}
// 	return table
// }
