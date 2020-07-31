package main

import (
	"fmt"
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

func addPanel(panelName string, x float64, y float64, width int, height int, win *pkg.Window, buttonAction1 func(event events.IEvent) error, buttonAction2 func(event events.IEvent) error) error {
	panel := widgets.NewPanel(panelName, x, y, width, height, nil)
	button := widgets.NewTextButton("button1", "my button1", 0, 0, 100, 100, nil, nil, nil)
	button.RegisterEventHandler(events.EventTypeButtonDown, buttonAction1)
	button2 := widgets.NewTextButton("button2", "my button2", 100, 0, 100, 100, nil, nil, nil)
	button2.RegisterEventHandler(events.EventTypeButtonDown, buttonAction2)
	panel.AddWidget(&button)
	panel.AddWidget(&button2)

	win.AddPanel(&panel)

	return nil
}

func main() {

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	log.SetLevel(log.DebugLevel)

	a := MyApp{}

	app := pkg.NewWindow(600, 600, "my title", true)
	addPanel("panel1", 100, 30, 200, 200, &app, a.ButtonAction1, a.ButtonAction2)
	addPanel("panel2", 300, 30, 200, 200, &app, a.ButtonAction1, a.ButtonAction2)

	// panel in panel
	panel3 := widgets.NewPanel("panel 3", 0, 300, 300, 300, &color.RGBA{0xff, 0xff, 0xff, 0xff})

	button := widgets.NewImageButton("image button 1", "./images/pressedbutton.png", "./images/nonpressedbutton.png", 0, 0)
	panel3.AddWidget(&button)

	cb := widgets.NewCheckBox("checkbox1", "./images/emptycheckbox.png", "./images/checkedcheckbox.png", 0, 100)
	cb.RegisterEventHandler(events.EventTypeButtonDown, a.CheckboxChanged)

	f := common.LoadFont("",16, color.RGBA{0xff,0xff,0xff,0xff})
	ti := widgets.NewTextInput("testinput1", 0, 150, 100, 20, &color.RGBA{0x55,0x55,0x55,0xff}, &f)

	panel3.AddWidget(&cb)
	panel3.AddWidget(&ti)

	button.RegisterEventHandler(events.EventTypeButtonDown, a.ButtonAction1)

	app.AddPanel(&panel3)

	go func(){
		for {
			time.Sleep(2 * time.Second)
			data, _ := ti.GetData()
			sData := data.(string)
			fmt.Printf("text is %s\n", sData)

			data2,_ := cb.GetData()
			cbData := data2.(bool)
			fmt.Printf("checkbox is %v\n", cbData)
		}
	}()


	app.MainLoop()

}
