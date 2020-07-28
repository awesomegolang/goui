package main

import (
	"github.com/kpfaulkner/goui/pkg"
	"github.com/kpfaulkner/goui/pkg/widgets"
)

func main() {
	app := pkg.NewWindow(600, 600, "my title")
	//button := widgets.NewImageButton("c:/temp/test.png",0,0,100,100)

	panel := widgets.NewPanel(0, 0, 200, 200)
	button := widgets.NewTextButton("my button", 0, 0, 100, 100, "", nil)
	panel.AddButton(&button)

	app.AddPanel(panel)
	app.MainLoop()

}