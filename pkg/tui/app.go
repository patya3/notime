package tui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/patya3/notime/pkg/tui/pages/main"
	"github.com/rivo/tview"
)

var app = tview.NewApplication()
var pagePrimitive = tview.NewPages()

func StartTui() {

	pagePrimitive.SetBackgroundColor(tcell.ColorDefault)

	pagePrimitive.
		AddPage("Main", mainPage.InitMainPage(app), true, true)

	if err := app.SetRoot(pagePrimitive, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
