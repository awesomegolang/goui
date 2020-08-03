package pkg

import (
	"fmt"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/kpfaulkner/goui/pkg/common"
	"github.com/kpfaulkner/goui/pkg/events"
	"github.com/kpfaulkner/goui/pkg/widgets"
	log "github.com/sirupsen/logrus"
	"image/color"
)

// Window used to define the UI window for the application.
// Currently will just cater for single window per app. This will be
// reviewed in the future.
type Window struct {
	width  int
	height int
	title  string

	// slice of panels. Should probably do as a map???
	// Then again, slice can be used for render order?
	panels []widgets.IPanel

	leftMouseButtonPressed  bool
	rightMouseButtonPressed bool

	haveMenuBar bool

	// These are other widgets/components that are listening to THiS widget. Ie we will broadcast to them!
	eventListeners map[int][]chan events.IEvent

	// incoming events to THIS widget (ie stuff we're listening to!)
	incomingEvents chan events.IEvent

	// widget that has focus...  I think that will do?
	FocusedWidget *widgets.IWidget
}

func NewWindow(width int, height int, title string, haveMenuBar bool) Window {
	w := Window{}
	w.height = height
	w.width = width
	w.title = title

	// panels are ordered. They are drawn from first to last.
	// So if we *have* to have a panel drawn last (eg, menu from a menu bar) then
	// one approach might be to create a panel (representing the menu) and it gets displayed at the end?
	// Mad idea..
	w.panels = []widgets.IPanel{}
	w.leftMouseButtonPressed = false
	w.rightMouseButtonPressed = false
	w.haveMenuBar = haveMenuBar

	w.eventListeners = make(map[int][]chan events.IEvent)
	w.incomingEvents = make(chan events.IEvent, 1000) // too much?

	if w.haveMenuBar {
		mb := *widgets.NewMenuBar("menubar", 0, 0, width, 30, &color.RGBA{0x71, 0x71, 0x71, 0xff})
		mb.AddMenuHeading("test")
		w.AddPanel(&mb)
	}
	return w
}

func (w *Window) AddEventListener(eventType int, ch chan events.IEvent) error {
	if _, ok := w.eventListeners[eventType]; ok {
		w.eventListeners[eventType] = append(w.eventListeners[eventType], ch)
	} else {
		w.eventListeners[eventType] = []chan events.IEvent{ch}
	}

	return nil
}

func (w *Window) RemoveEventListener(eventType int, ch chan events.IEvent) error {
	if _, ok := w.eventListeners[eventType]; ok {
		for i := range w.eventListeners[eventType] {
			if w.eventListeners[eventType][i] == ch {
				w.eventListeners[eventType] = append(w.eventListeners[eventType][:i], w.eventListeners[eventType][i+1:]...)
				break
			}
		}
	}
	return nil
}

func (w *Window) GetEventListenerChannel() chan events.IEvent {
	return w.incomingEvents
}

// Emit event for  all listeners to receive
func (w *Window) EmitEvent(event events.IEvent) error {
	if _, ok := w.eventListeners[event.EventType()]; ok {
		for _, handler := range w.eventListeners[event.EventType()] {
			go func(handler chan events.IEvent) {
				handler <- event
			}(handler)
		}
	}

	return nil
}

func (w *Window) AddPanel(panel widgets.IPanel) error {
	panel.SetTopLevel(true)
	w.panels = append(w.panels, panel)
	return nil
}

// FindWidgetForInput
// Need to make recursive for panels in panels etc... but just leave pretty linear for now.
func (w *Window) FindWidgetForInput(x float64, y float64) (*widgets.IWidget, error) {

	// all things are panels at this level.
	for _, panel := range w.panels {
		if panel.ContainsCoords(x, y) {

			for _, subPanel := range panel.ListPanels() {
				if subPanel.ContainsCoords(x, y) {
					for _, widget := range subPanel.ListWidgets() {
						if widget.ContainsCoords(x, y) {
							// have match
							w.FocusedWidget = &widget
							return &widget,nil
						}
					}
				}
			}

			// check widgets in panel.
			for _, widget := range panel.ListWidgets() {
				//px,py := panel.GetCoords()
				//xx := x - px
				//yy := y - py
				if widget.ContainsCoords(x,y) {
					// have match
					w.FocusedWidget = &widget
					return &widget,nil
				}
			}
		}
	}

	//return nil, errors.New("Unable to panel/widget that was clicked on....  impossible!!!")
	return nil,nil
}

/////////////////////// EBiten specifics below... /////////////////////////////////////////////
func (w *Window) Update(screen *ebiten.Image) error {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		w.leftMouseButtonPressed = true
		x, y := ebiten.CursorPosition()
		usedWidget, err := w.FindWidgetForInput(float64(x), float64(y))
		if err != nil {
			log.Errorf("Unable to find widget for click!! %s", err.Error())
		}

		if usedWidget != nil {
			me := events.NewMouseEvent(fmt.Sprintf("widget %s button down", (*usedWidget).GetID()), x,y, events.EventTypeButtonDown)
			//w.EmitEvent(me)
			//(*usedWidget).HandleEvent(me)
			(*usedWidget).HandleEvent(me)
		}
	} else {
		if w.leftMouseButtonPressed {
			w.leftMouseButtonPressed = false

			x, y := ebiten.CursorPosition()
			usedWidget, err := w.FindWidgetForInput(float64(x), float64(y))
			if err != nil {
				log.Errorf("Unable to find widget for click!! %s", err.Error())
			}

			if usedWidget != nil {
				me := events.NewMouseEvent(fmt.Sprintf("widget %s button up", (*usedWidget).GetID()), x,y, events.EventTypeButtonUp)
				//w.EmitEvent(me)
				//(*usedWidget).HandleEvent(me)
				(*usedWidget).HandleEvent(me)
			}

			// it *WAS* pressed previous frame... but isn't now... this means released!!!
			//x, y := ebiten.CursorPosition()
			//me := events.NewMouseEvent(x, y, events.EventTypeButtonUp)
			//w.EmitEvent(me)
		}
	}

	inp := ebiten.InputChars()
	if len(inp) > 0 {
		// create keyboard event

		//w.EmitEvent(ke)
		x, y := ebiten.CursorPosition()
		usedWidget, err := w.FindWidgetForInput(float64(x), float64(y))
		if err != nil {
			log.Errorf("Unable to find widget for click!! %s", err.Error())
		}

		if usedWidget != nil {
			//ke := events.NewMouseEvent(x,y, events.EventTypeKeyboard)
			ke := events.NewKeyboardEvent(ebiten.Key(inp[0])) // only send first one?
			//w.EmitEvent(me)
			(*usedWidget).HandleEvent(ke)
		}
	}

	// If the backspace key is pressed, remove one character.
	if repeatingKeyPressed(ebiten.KeyBackspace) {
		x, y := ebiten.CursorPosition()
		usedWidget, err := w.FindWidgetForInput(float64(x), float64(y))
		if err != nil {
			log.Errorf("Unable to find widget for click!! %s", err.Error())
		}

		if usedWidget != nil {
			//ke := events.NewMouseEvent(x,y, events.EventTypeKeyboard)
			ke := events.NewKeyboardEvent(ebiten.KeyBackspace)
			//w.EmitEvent(me)
			//(*usedWidget).HandleEvent(ke)
			(*usedWidget).HandleEvent(ke)
		}
	}

	return nil
}

func (w *Window) HandleButtonUpEvent(event events.MouseEvent) error {
	log.Debugf("button up %f %f", event.X, event.Y)
	for _, panel := range w.panels {
		panel.HandleEvent(event)
	}

	return nil
}

func (w *Window) HandleButtonDownEvent(event events.MouseEvent) error {
	log.Debugf("button down %f %f", event.X, event.Y)

	// loop through panels and find a target!
	for _, panel := range w.panels {
		panel.HandleEvent(event)
	}

	return nil
}

func (w *Window) HandleKeyboardEvent(event events.KeyboardEvent) error {

	// loop through panels and find a target!
	for _, panel := range w.panels {
		panel.HandleEvent(event)
	}
	return nil
}

func (w *Window) HandleEvent(event events.IEvent) error {
	//log.Debugf("Window handled event %v", event)

	switch event.EventType() {
	case events.EventTypeButtonUp:
		{
			err := w.HandleButtonUpEvent(event.(events.MouseEvent))
			return err
		}

	case events.EventTypeButtonDown:
		{
			err := w.HandleButtonDownEvent(event.(events.MouseEvent))
			return err
		}

	case events.EventTypeKeyboard:
		{
			err := w.HandleKeyboardEvent(event.(events.KeyboardEvent))
			return err
		}
	}

	return nil
}

func (w *Window) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x11, 0x11, 0x11, 0xff})

	for _, panel := range w.panels {
		panel.Draw(screen)
	}

	x, y := ebiten.CursorPosition()
	defaultFontInfo := common.LoadFont("", 16, color.RGBA{0xff, 0xff, 0xff, 0xff})
	text.Draw(screen, fmt.Sprintf("%d %d", x,y), defaultFontInfo.UIFont, 00, 500, color.White)

}

func (w *Window) Layout(outsideWidth, outsideHeight int) (int, int) {
	return w.width, w.height
}

func (w *Window) MainLoop() error {

	ebiten.SetWindowSize(w.width, w.height)
	ebiten.SetWindowTitle(w.title)
	if err := ebiten.RunGame(w); err != nil {
		log.Fatal(err)
	}

	return nil
}

// repeatingKeyPressed return true when key is pressed considering the repeat state.
func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 30
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}
