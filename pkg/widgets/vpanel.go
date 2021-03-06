package widgets

import (
	"image/color"
)

// VPanel stacks internal objects vertically.... starting at top and going down... so umm reverse panel :)
type VPanel struct {
	Panel

	// Y offset for next widget to be added.
	YLoc float64

	resizable bool
}

func NewVPanel(ID string, colour *color.RGBA) *VPanel {
	p := VPanel{}
	p.Panel = *NewPanel(ID, colour, nil)
	p.DynamicSize = true
	return &p
}

func NewVPanelWithSize(ID string, width int, height int, colour *color.RGBA) *VPanel {
	p := NewVPanel(ID, colour)
	p.SetSize(width, height)
	p.DynamicSize = false
	return p
}

func (p *VPanel) ClearWidgets() error {
	p.widgets = []IWidget{}
	p.YLoc = 0
	return nil
}

// AddWidget adds a widget to the panel, but each widget is below the next one.
func (p *VPanel) AddWidget(w IWidget) error {
	w.AddParentPanel(p)

	// find X,Y for widget...
	w.SetXY(p.X, p.YLoc)
	width, height := w.GetSize()

	if p.DynamicSize {
		// grow panel height if widget is taller.
		if p.YLoc+height > float64(p.Height) {
			p.Height = int(p.YLoc + height)
			p.SetSize(p.Width, p.Height)
		}

		if width > float64(p.Width) {
			p.Width = int(width)
			p.SetSize(p.Width, p.Height)
		}
	}
	p.YLoc += height

	p.widgets = append(p.widgets, w)
	return nil
}
