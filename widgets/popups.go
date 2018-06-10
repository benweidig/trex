package widgets

import (
	"github.com/benweidig/trex/input"
	"github.com/benweidig/trex/ui"
	"github.com/rivo/tview"
)

// NewFormatterPopup builds a new tview.Primitive for the formatter chooser
func NewFormatterPopup(selectedFn func(fileType input.FileType)) tview.Primitive {
	l := ui.NewList()
	l.SetRect(0, 0, 10, 4)
	items := []ui.ListItem{
		ui.NewSimpleListItem("  JSON  "),
		ui.NewSimpleListItem("  YAML  "),
	}
	l.SetItems(items, false)
	l.SetBorder(true)
	l.SetTitle("Output")

	l.SetSelectedFn(func(idx int, item ui.ListItem) {
		switch idx {
		case 0:
			selectedFn(input.FileTypeJSON)

		case 1:
			selectedFn(input.FileTypeYAML)
		}
	})

	return ui.NewPopup(l)
}

const helpPopupText = `(shift + ←/→) Resize
        (tab) Switch focus (tree/output)
          (f) Choose Formatter
          (c) Copy currently selected node
          (?) Display help

Navigate with arrow keys / vim-keys`

// NewHelpPopup builds a tview.Primitive for displaying some help text
func NewHelpPopup() tview.Primitive {
	t := tview.NewTextView()
	t.SetBorder(true)
	t.SetTitle(" Help ")
	t.SetRect(0, 0, 47, 9)
	t.SetBorderPadding(1, 1, 1, 1)
	t.SetText(helpPopupText)

	return ui.NewPopup(t)
}
