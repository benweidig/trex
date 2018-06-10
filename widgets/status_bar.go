package widgets

import (
	"strings"

	"github.com/benweidig/trex/input"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// StatusBar is a  smal wrapper around tview.TextView representing a simple 1-lined status bar
type StatusBar struct {
	textView *tview.TextView
	path     string
	fileType input.FileType
}

// NewStatusBar creates a new StatusBar
func NewStatusBar() *StatusBar {
	b := &StatusBar{
		textView: tview.NewTextView(),
	}

	b.textView.SetTextColor(tview.Styles.PrimitiveBackgroundColor)
	b.textView.SetBackgroundColor(tview.Styles.PrimaryTextColor)
	return b
}

// SetContent updates the displayed content of the bar
func (b *StatusBar) SetContent(path string, fileType input.FileType) *StatusBar {
	b.path = path
	b.fileType = fileType
	return b
}

// SetPath only updates the path, filetype remains
func (b *StatusBar) SetPath(path string) *StatusBar {
	b.SetContent(path, b.fileType)
	return b
}

// SetFileType only updates the filetype, path remains
func (b *StatusBar) SetFileType(fileType input.FileType) *StatusBar {
	b.SetContent(b.path, fileType)
	return b
}

// Draw implements tview.Primitive
func (b *StatusBar) Draw(screen tcell.Screen) {
	_, _, width, _ := b.textView.GetInnerRect()

	fileType := string(b.fileType)
	actualWidth := width - 2
	paddingWidth := actualWidth - len(b.path) - len(fileType)
	padding := strings.Repeat(" ", paddingWidth)

	b.textView.SetText(" " + b.path + padding + fileType)

	b.textView.Draw(screen)
}

// GetRect implements tview.Primitive
func (b *StatusBar) GetRect() (int, int, int, int) {
	return b.textView.GetRect()
}

// SetRect implements tview.Primitive
func (b *StatusBar) SetRect(x, y, width, height int) {
	b.textView.SetRect(x, y, width, height)
}

// InputHandler implements tview.Primitive
func (b *StatusBar) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return b.textView.InputHandler()
}

// Focus implements tview.Primitive
func (b *StatusBar) Focus(delegate func(p tview.Primitive)) {
	b.textView.Focus(delegate)
}

// Blur implements tview.Primitive
func (b *StatusBar) Blur() {
	b.textView.Blur()
}

// GetFocusable implements tview.Primitive
func (b *StatusBar) GetFocusable() tview.Focusable {
	return b.textView.GetFocusable()
}
