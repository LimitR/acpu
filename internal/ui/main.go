package ui

import (
	"strconv"

	"urinal/internal/models"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const asciiUrinal = `██                                  ██
██             ████        █        ██
██             ████▄▄▄██████        ██
██             ████▌ ▐     █        ██
██                 ▌ ▐              ██
██                 ▌ ▐              ██
██                 ▌ ▐              ██
██                 ▌ ▐              ██
██      ███████████████████████     ██
██      ██░░░░░░░░░░░░░░░░░░░██     ██
██      ██░░███████████████░░██     ██
██      ██░░█▓▓▓▓▓▓▓▓▓▓▓▓▓█░░██     ██
██      ██░░█▓▓▓▓▓▓▓▓▓▓▓▓▓█░░██     ██
██      ██░░█▓▓▓▓▓▓▓▓▓▓▓▓▓█░░██     ██
██      ██░░█▓▓▓▓▓▓▓▓▓▓▓▓▓█░░██     ██
██      ██░░█▓▓▓▓▓▓▓▓▓▓▓▓▓█░░██     ██
██      ██░░█▓▓▓▓▓▓▓▓▓▓▓▓▓█░░██     ██
██      ██░░█▓▓▓▓▓▓▓▓▓▓▓▓▓█░░██     ██
██      ██░░█▓▓▓▓▓▓▓▓▓▓▓▓▓█░░██     ██
██      ██░░█▓▓▓▓▓▓▓▓▓▓▓▓▓█░░██     ██
██      ██░░░██▓▓▓▓▓▓▓▓▓██░░░██     ██
██      ██░░░░░██▓▓▓▓▓██░░░░░██     ██
██       ██░░░░░░█████░░░░░░██      ██
██        ██░░░░░░░░░░░░░░░██       ██
██         ██░░░░░░░░░░░░░██        ██
██          ██░░░░░░░░░░░██         ██
██            ███████████           ██`

func RunUI() {
	app := tview.NewApplication()
	arrCheckboxIndex := make(map[int]int, 10)
	pages := tview.NewPages()
	t := models.NewToilet()
	changeNumber := 0
	form := tview.NewForm().SetHorizontal(true)
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyESC || key == tcell.KeyENQ {
				app.Stop()
			}
			form.Clear(true)
			pages.SwitchToPage("page-1")
		}).SetTextAlign(20)
	inp := tview.NewInputField().SetLabel("Enter a quantity urinal: ").
		SetFieldWidth(15).
		SetAcceptanceFunc(tview.InputFieldInteger).
		SetDoneFunc(func(key tcell.Key) {
			t.AddUrinal(changeNumber)
			for i := 0; i < changeNumber; i++ {
				func(inp int) {
					form.AddCheckbox("[green]"+strconv.Itoa(inp+1), false, func(checked bool) {
						if checked {
							arrCheckboxIndex[inp] = inp
						} else {
							delete(arrCheckboxIndex, inp)
						}
					})
				}(i)
			}
			form.AddButton("Search optimal place", func() {
				for _, v := range arrCheckboxIndex {
					err := t.TakePlace(v)
					if err != nil {
						textView.SetText(err.Error())
						pages.SwitchToPage("page-3")
					}
				}
				optimalIndex, err := t.GetOptimalPlace()
				needIndex := ""
				resultTextNumber := " "
				resultText := ""
				if err != nil {
					resultText = err.Error()
				}
				for i := 0; i < len(t.ArrayUrinal); i++ {
					if optimalIndex == i {
						resultTextNumber += " [green]" + strconv.Itoa(i+1) + " "
						needIndex = "[green]" + strconv.Itoa(i+1)
						resultText += "[green]🚽 "
						continue
					}
					if t.ArrayUrinal[i] == len(t.ArrayUrinal) {
						resultTextNumber += " [red]" + strconv.Itoa(i+1) + " "
						resultText += "[red]🧍‍♂️ "
						continue
					}
					resultTextNumber += " [white]" + strconv.Itoa(i+1) + " "
					resultText += "[white]🚽 "
				}
				textView.SetText(resultTextNumber + "\n" + resultText + "\n[white]Your urinal is under the number " + needIndex + "\n\n\n" + asciiUrinal).SetTextAlign(20)
				pages.SwitchToPage("page-3")
			})
			pages.SwitchToPage("page-2")
		}).
		SetChangedFunc(func(text string) {
			changeNumber, _ = strconv.Atoi(text)
		})
	pages.AddPage("page-3", textView, true, true)
	pages.AddPage("page-2", form, true, true)
	pages.AddPage("page-1", inp, true, true)

	box := tview.NewFlex().
		AddItem(pages, 0, 1, true)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyCtrlQ {
			app.Stop()
		}
		return event
	})
	if err := app.SetRoot(box, true).EnableMouse(true).SetFocus(box).Run(); err != nil {
		panic(err)
	}
}
