package widgets

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Output is a simple tview.TextView wrapper to simplify the exposed interface
type Output struct {
	textView *tview.TextView
}

// NewOutput creates new simple Output widget
func NewOutput(monochrome bool) *Output {
	o := &Output{
		textView: tview.NewTextView(),
	}
	o.textView.SetDynamicColors(!monochrome)
	o.textView.SetScrollable(true)
	o.textView.SetBorder(true)

	return o
}

// Draw implements tview.Primitive
func (o *Output) Draw(screen tcell.Screen) {
	o.textView.Draw(screen)
}

// GetRect implements tview.Primitive
func (o *Output) GetRect() (int, int, int, int) {
	return o.textView.GetRect()
}

// SetRect implements tview.Primitive
func (o *Output) SetRect(x int, y int, width int, height int) {
	o.textView.SetRect(x, y, width, height)
}

// InputHandler implements tview.Primitive, restriction to single line navigation only
func (o *Output) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return o.textView.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if o.textView.HasFocus() == false {
			return
		}

		if event.Key() == tcell.KeyDown || event.Key() == tcell.KeyUp {
			o.textView.InputHandler()(event, setFocus)
			return
		}

		if event.Key() == tcell.KeyRune {
			if event.Rune() == 'j' || event.Rune() == 'k' {
				o.textView.InputHandler()(event, setFocus)
				return
			}
		}
	})
}

// Focus implements tview.Primitive
func (o *Output) Focus(delegate func(p tview.Primitive)) {
	o.textView.Focus(delegate)
}

// Blur implements tview.Primitive
func (o *Output) Blur() {
	o.textView.Blur()
}

// GetFocusable implements tview.Primitive
func (o *Output) GetFocusable() tview.Focusable {
	return o.textView.GetFocusable()
}

// SetText sets the text and scrolls to the beginning
func (o *Output) SetText(text string) *Output {
	o.textView.SetText(text)
	o.textView.ScrollToBeginning()
	return o
}
