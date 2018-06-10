package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Popup is a simple helper for centering a primitive horizontally/vertically.
// It contains the actual primitive and delegates all related calls to it.
type Popup struct {
	flex    *tview.Flex
	content tview.Primitive
}

// NewPopup creates a new horizontal-/vertical-centered Popup.
func NewPopup(p tview.Primitive) *Popup {
	_, _, width, height := p.GetRect()
	popup := &Popup{
		flex: tview.NewFlex().
			AddItem(nil, 0, 1, false).
			AddItem(tview.NewFlex().
				SetDirection(tview.FlexRow).
				AddItem(nil, 0, 1, false).
				AddItem(p, height, 1, false).
				AddItem(nil, 0, 1, false), width, 1, false).
			AddItem(nil, 0, 1, false),
	}
	popup.content = p
	return popup
}

// Draw implements tview.Primitive
func (p *Popup) Draw(screen tcell.Screen) {
	p.flex.Draw(screen)
}

// GetRect implements tview.Primitive
func (p *Popup) GetRect() (int, int, int, int) {
	return p.flex.GetRect()
}

// SetRect implements tview.Primitive
func (p *Popup) SetRect(x, y, width, height int) {
	p.flex.SetRect(x, y, width, height)
}

// InputHandler implements tview.Primitive
func (p *Popup) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return p.content.InputHandler()
}

// Focus implements tview.Primitive
func (p *Popup) Focus(delegate func(p tview.Primitive)) {
	p.content.Focus(delegate)
}

// Blur implements tview.Primitive
func (p *Popup) Blur() {
	p.content.Blur()
}

// GetFocusable implements tview.Primitive
func (p *Popup) GetFocusable() tview.Focusable {
	return p.content.GetFocusable()
}
