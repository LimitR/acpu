package ui

import (
	"strconv"

	"urinal/internal/models"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func RunUI() {
	app := tview.NewApplication()
	pages := tview.NewPages()
	t := models.NewToilet()
	changeNumber := 0
	form := tview.NewForm().SetHorizontal(true)
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true)
	inp := tview.NewInputField().SetLabel("Enter a quantity urinal: ").
		SetFieldWidth(15).
		SetAcceptanceFunc(tview.InputFieldInteger).
		SetDoneFunc(func(key tcell.Key) {
			t.AddUrinal(changeNumber)
			for i := 0; i < changeNumber; i++ {
				func(inp int) {
					form.AddCheckbox(strconv.Itoa(inp+1), false, func(checked bool) {
						err := t.TakePlace(inp)
						if err != nil {
							textView.SetText(err.Error())
							pages.SwitchToPage("page-3")
						}
					})
				}(i)
			}
			form.AddButton("Search optimal place", func() {
				optimalIndex, err := t.GetOptimalPlace()
				needIndex := ""
				resultText := ""
				if err != nil {
					resultText = err.Error()
				}
				for i := 0; i < len(t.ArrayUrinal); i++ {
					if optimalIndex == i {
						needIndex = "[green]" + strconv.Itoa(i+1)
						resultText += "[green]" + strconv.Itoa(i+1) + "█[red] "
						continue
					}
					resultText += "[red]" + strconv.Itoa(i+1) + "█ "
				}
				textView.SetText(resultText + "\nYour urinal is under the number " + needIndex)
				pages.SwitchToPage("page-3")
			})
			pages.SwitchToPage("page-2")
		}).
		SetChangedFunc(func(text string) {
			changeNumber, _ = strconv.Atoi(text)
		})
	pages.AddPage("page-3", textView, true, false)
	pages.AddPage("page-2", form, true, false)
	pages.AddPage("page-1", inp, true, true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlQ {
			app.Stop()
		}
		return event
	})
	if err := app.SetRoot(pages, true).SetFocus(pages).Run(); err != nil {
		panic(err)
	}
}
