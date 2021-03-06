package widgets

import (
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/kpfaulkner/goui/pkg/common"
	"github.com/kpfaulkner/goui/pkg/events"
	log "github.com/sirupsen/logrus"
	"golang.org/x/image/font"
	"image/color"
	_ "image/png"
)

var defaultFontInfo common.Font

type TextInput struct {
	BaseWidget

	text             string
	backgroundColour color.RGBA
	fontInfo         common.Font
	uiFont           font.Face

	// just for cursor.
	counter int

	// vertical position for text
	vertPos int
}

func init() {
	defaultFontInfo = common.LoadFont("", 16, color.RGBA{0xff, 0xff, 0xff, 0xff})
}

func NewTextInput(ID string, width int, height int, backgroundColour *color.RGBA, fontInfo *common.Font, handler func(event events.IEvent) error) *TextInput {
	t := TextInput{}
	t.BaseWidget = *NewBaseWidget(ID, width, height, handler)
	t.text = ""
	t.stateChangedSinceLastDraw = true
	t.counter = 0

	if backgroundColour != nil {
		t.backgroundColour = *backgroundColour
	} else {
		t.backgroundColour = color.RGBA{0, 0xff, 0, 0xff}
	}

	if fontInfo != nil {
		t.fontInfo = *fontInfo
	} else {
		t.fontInfo = defaultFontInfo
	}

	// vert pos is where does text go within button. Assuming we want it centred (for now)
	// Need to just find something visually appealing.
	t.vertPos = (height - (height-int(t.fontInfo.SizeInPixels))/2) - 2

	return &t
}

func (t *TextInput) HandleEvent(event events.IEvent) error {

	eventType := event.EventType()
	switch eventType {
	case events.EventTypeButtonDown:
		{
			mouseEvent := event.(events.MouseEvent)

			// check click is in button boundary.
			if t.ContainsCoords(mouseEvent.X, mouseEvent.Y) {
				t.hasFocus = true
				t.stateChangedSinceLastDraw = true
				// then do application specific stuff!!

			} else {
				t.hasFocus = false
			}
		}
	case events.EventTypeKeyboard:
		{
			// check if has focus....  if so, can potentially add to string?
			if t.hasFocus {
				keyboardEvent := event.(events.KeyboardEvent)

				if keyboardEvent.Character != ebiten.KeyBackspace {
					t.text = t.text + string(keyboardEvent.Character)
				} else {
					// back space one.
					l := len(t.text)
					if l > 0 {
						t.text = t.text[0 : l-1]
					}
				}
				t.stateChangedSinceLastDraw = true
			}
		}
	case events.EventTypeSetText:
		{
			setTextEvent := event.(events.SetTextEvent)

			t.text = setTextEvent.Text
			t.stateChangedSinceLastDraw = true
		}
	}

	return nil
}

func (t *TextInput) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(t.X, t.Y)

	if t.stateChangedSinceLastDraw {
		log.Debugf("textinput text %s", t.text)
		// how often do we update this?
		emptyImage, _ := ebiten.NewImage(t.Width, t.Height, ebiten.FilterDefault)
		_ = emptyImage.Fill(t.backgroundColour)
		t.rectImage = emptyImage

		ebitenutil.DrawLine(t.rectImage, 0, 0, float64(t.Width), 0, color.Black)
		ebitenutil.DrawLine(t.rectImage, float64(t.Width), 0, float64(t.Width), float64(t.Height), color.Black)
		ebitenutil.DrawLine(t.rectImage, float64(t.Width), float64(t.Height), 0, float64(t.Height), color.Black)
		ebitenutil.DrawLine(t.rectImage, 0, float64(t.Height), 0, 0, color.Black)

		txt := t.text
		txt += "|"

		text.Draw(t.rectImage, txt, t.fontInfo.UIFont, 0, t.vertPos, t.fontInfo.Colour)
		t.stateChangedSinceLastDraw = false
	}

	// if state changed since last draw, recreate colour etc.
	_ = screen.DrawImage(t.rectImage, op)

	return nil
}

func (t *TextInput) GetData() (interface{}, error) {
	return t.text, nil
}
