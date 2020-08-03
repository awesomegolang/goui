package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/kpfaulkner/goui/pkg"
	"github.com/kpfaulkner/goui/pkg/common"
	"github.com/kpfaulkner/goui/pkg/events"
	"github.com/kpfaulkner/goui/pkg/widgets"
	log "github.com/sirupsen/logrus"
	"image/color"
	"net/http"
	_ "net/http/pprof"
	"time"
)

type MyApp struct {
	mytext string
}

func (m *MyApp) ButtonAction1(event events.IEvent) error {
	log.Debugf("My button1 action!!!")
	return nil
}

func (m *MyApp) ButtonAction2(event events.IEvent) error {
	log.Debugf("My button2 action!!!")
	return nil
}

func (m *MyApp) CheckboxChanged(event events.IEvent) error {
	log.Debugf("checkbox changed!!!")
	return nil
}

func (m *MyApp) HandleTextInput(event events.IEvent) error {
	log.Debugf("text input changed.!!!")

	kbEvent := event.(events.KeyboardEvent)
	m.mytext += string(kbEvent.Character)
	return nil
}

func addPanel(panelName string, width int, height int, win *pkg.Window, buttonAction1 func(event events.IEvent) error, buttonAction2 func(event events.IEvent) error) error {
	panel := widgets.NewPanel(panelName, width, height, nil)

	button := widgets.NewTextButton("button1", "my button1",  100, 100, nil, nil, nil)
	panel.AddWidget(button)

	button2 := widgets.NewTextButton("button2", "my button2",  100, 100, nil, nil, nil)
	panel.AddWidget(button2)

	win.AddPanel(panel)

	return nil
}

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.SetLevel(log.DebugLevel)

	a := MyApp{}

	app := pkg.NewWindow(600, 600, "my title", false)
	addPanel("panel1",  200, 200, &app, a.ButtonAction1, a.ButtonAction2)

	addPanel("panel2",  200, 200, &app, a.ButtonAction1, a.ButtonAction2)

	panel := widgets.NewPanel("panel3",  200, 200, nil)
	button := widgets.NewImageButton("image button 1", "./images/pressedbutton.png", "./images/nonpressedbutton.png")
	panel.AddWidget(button)

	cb := widgets.NewCheckBox("checkbox1", "./images/emptycheckbox.png", "./images/checkedcheckbox.png", 0, 100)
	panel.AddWidget(cb)

	//cb.RegisterEventHandler(events.EventTypeButtonDown, a.CheckboxChanged)

	f := common.LoadFont("", 16, color.RGBA{0xff, 0xff, 0xff, 0xff})
	ti := widgets.NewTextInput("testinput1",100, 20, &color.RGBA{0x55, 0x55, 0x55, 0xff}, &f)
	panel.AddWidget(ti)

	app.AddPanel(panel)

	go func() {
		for {
			time.Sleep(2 * time.Second)
			data, _ := ti.GetData()
			sData := data.(string)
			fmt.Printf("text is %s\n", sData)

			data2, _ := cb.GetData()
			cbData := data2.(bool)
			fmt.Printf("checkbox is %v\n", cbData)
		}
	}()

	ebiten.SetRunnableInBackground(true)
	app.MainLoop()

}
