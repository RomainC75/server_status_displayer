package table

import (
	"fmt"
	"time"
	"tviewTest/tester"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func draw() {
	app := tview.NewApplication()
	table := tview.NewTable()

	go func() {
		i := 0
		for {
			app.QueueUpdateDraw(func() {
				table.Clear()
				for row := 0; row < 3; row++ {
					for column := 0; column < 3; column++ {
						color := tview.AlignRight

						var txt string
						switch column {
						case 0:
							txt = data[row].Date.String()
						case 1:
							txt = data[row].Region
						case 2:
							txt = fmt.Sprintf(" xx%.1fxx ", data[row].Price)
							// txt = strconv.Itoa(i)
						}

						table.SetCell(row, column,
							tview.NewTableCell(txt).
								SetTextColor(tview.Styles.PrimaryTextColor).
								SetAlign(color))
					}
				}
				i++
			})
			time.Sleep(1 * time.Second)
		}
	}()

	go func() {
		j := 0
		for {

			data[2].Price = float32(j)
			j++
			time.Sleep(time.Second * 2)
		}
	}()

	if err := app.SetRoot(table, true).Run(); err != nil {
		panic(err)
	}
}

func Merger(channels []chan tester.TestingResult) {

}

func display(data []Data) *tview.Table {
	table := tview.NewTable().SetFixed(1, 1)
	for row := 0; row < len(data); row++ {
		for column := 0; column < 3; column++ {
			color := tcell.ColorWhite
			if row == 0 {
				color = tcell.ColorYellow
			} else if column == 0 {
				color = tcell.ColorDarkCyan
			}
			align := tview.AlignLeft
			if row == 0 {
				align = tview.AlignCenter
			} else if column == 0 || column >= 4 {
				align = tview.AlignRight
			}

			var txt string
			switch column {
			case 0:
				txt = data[row].Date.String()
			case 1:
				txt = data[row].Region
			case 2:
				txt = fmt.Sprintf("%f", data[row].Price)
			}
			table.SetCell(row, column, &tview.TableCell{
				Text:  txt,
				Color: color,
				Align: align,
			})
		}
	}
	return table
}
