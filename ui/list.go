package ui

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// ListItem represents the minimal item that a List can use.
type ListItem interface {
	Label() string
}

// SimpleListItem is a helper for easier ListItem use.
type SimpleListItem struct {
	label string
}

// NewSimpleListItem returns a new SimpleListItem
func NewSimpleListItem(label string) ListItem {
	return &SimpleListItem{
		label: label,
	}
}

// Label implements the interface ListItem
func (i SimpleListItem) Label() string {
	return i.label
}

// List is a simple selectable List with callback function
type List struct {
	*tview.Box

	items []ListItem

	currentIdx  int
	currentItem ListItem

	overflow               bool
	textColor              tcell.Color
	currentTextColor       tcell.Color
	currentBackgroundColor tcell.Color

	// changedFn is called when a new item is navigated to.
	changedFn func(index int, listItem ListItem)

	// selectedFn is called when the user interacts with the current item (space, enter)
	selectedFn func(index int, listItem ListItem)

	// doneFn is called on Escape
	doneFn func()
}

// NewList build a List with default styles
func NewList() *List {
	return &List{
		Box:                    tview.NewBox(),
		textColor:              tview.Styles.PrimaryTextColor,
		currentTextColor:       tview.Styles.PrimitiveBackgroundColor,
		currentBackgroundColor: tview.Styles.PrimaryTextColor,
	}
}

// SetOverflow sets if the list should overflow after reaching the start/end.
func (l *List) SetOverflow(overflow bool) *List {
	l.overflow = overflow
	return l
}

// SetTextColor sets the color of non-highlighted text.
func (l *List) SetTextColor(color tcell.Color) *List {
	l.textColor = color
	return l
}

// SetCurrentTextColor sets the color of the currently highlighted item text.
func (l *List) SetCurrentTextColor(color tcell.Color) *List {
	l.currentTextColor = color
	return l
}

// SetCurrentBackgroundColor set the color of the currently highlighted item background.
func (l *List) SetCurrentBackgroundColor(color tcell.Color) *List {
	l.currentBackgroundColor = color
	return l
}

// SetChangedFn sets the callback function after the position of the current item has changed.
func (l *List) SetChangedFn(fn func(index int, item ListItem)) *List {
	l.changedFn = fn
	return l
}

// SetSelectedFn sets the callback function after interaction (enter/space) with the current item.
func (l *List) SetSelectedFn(fn func(index int, item ListItem)) *List {
	l.selectedFn = fn
	return l
}

// SetDoneFn sets the callback function if Escape is pressed
func (l *List) SetDoneFn(fn func()) *List {
	l.doneFn = fn
	return l
}

// ClearItems removes all list items.
func (l *List) ClearItems() *List {
	if len(l.items) > 0 {
		l.items = l.items[:0]
	}
	return l
}

// AddItem appends a new ListItem to the list.
func (l *List) AddItem(item ListItem) *List {
	l.items = append(l.items, item)
	return l
}

// GetItemCount returns the items count.
func (l *List) GetItemCount() int {
	return len(l.items)
}

// GetItems return the items.
func (l *List) GetItems() []ListItem {
	return l.items
}

// SetItems replaced the items.
// You can choose to keep the current position, with fallback to 0 if no longer valid.
func (l *List) SetItems(items []ListItem, keepPosition bool) *List {
	l.items = items
	if keepPosition == false || l.currentIdx > len(l.items)-1 {
		l.currentIdx = 0
	}
	l.currentItem = l.items[l.currentIdx]

	l.TriggerChanged()
	return l
}

// GetCurrentIdx returns the index of the current ListItem.
func (l *List) GetCurrentIdx() int {
	return l.currentIdx
}

// GetCurrentItem returns the current ListItem
func (l *List) GetCurrentItem() ListItem {
	return l.currentItem
}

// SetCurrentItem sets the current ListItem to a specific index.
func (l *List) SetCurrentItem(index int) *List {
	if l.currentIdx == index {
		return l
	}

	if l.overflow {
		if index < 0 {
			l.currentIdx = len(l.items) - 1
		} else if index >= len(l.items) {
			l.currentIdx = 0
		} else {
			l.currentIdx = index
		}
	} else {
		if index < 0 {
			l.currentIdx = 0
		} else if index >= len(l.items) {
			l.currentIdx = len(l.items) - 1
		} else {
			l.currentIdx = index
		}
	}

	item := l.items[l.currentIdx]
	l.currentItem = item
	l.TriggerChanged()
	return l
}

// TriggerChanged trigger the changedFn callback if set.
func (l *List) TriggerChanged() *List {
	if l.changedFn != nil {
		l.changedFn(l.currentIdx, l.currentItem)
	}
	return l
}

// Draw implements tview.Primitive
func (l *List) Draw(screen tcell.Screen) {
	l.Box.Draw(screen)

	// Determine the dimensions.
	x, y, width, height := l.GetInnerRect()
	bottomLimit := y + height

	// We want to keep the current selection in view. What is our offset?
	var offset int
	if l.currentIdx >= height {
		offset = l.currentIdx + 1 - height
	}

	// Draw the list items.
	for index, item := range l.items {
		if index < offset {
			continue
		}

		if y >= bottomLimit {
			break
		}

		// Main text.
		tview.Print(screen, item.Label(), x, y, width, tview.AlignLeft, l.textColor)

		// Background color of selected text.
		if index == l.currentIdx {
			textWidth := width
			for bx := 0; bx < textWidth && bx < width; bx++ {
				m, c, style, _ := screen.GetContent(x+bx, y)
				fg, _, _ := style.Decompose()
				if fg == l.textColor {
					fg = l.currentTextColor
				}
				style = style.Background(l.currentBackgroundColor).Foreground(fg)
				screen.SetContent(x+bx, y, m, c, style)
			}
		}

		y++

		if y >= bottomLimit {
			break
		}
	}
}

// InputHandler implements tview.Primitive
func (l *List) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
	return l.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
		if l.HasFocus() == false {
			return
		}

		newIndex := l.currentIdx

		switch key := event.Key(); key {
		case tcell.KeyLeft:
			newIndex--

		case tcell.KeyDown:
			newIndex++

		case tcell.KeyRight:
			newIndex++

		case tcell.KeyUp:
			newIndex--

		case tcell.KeyHome:
			newIndex = 0

		case tcell.KeyEnd:
			newIndex = len(l.items) - 1

		case tcell.KeyPgDn:
			newIndex += 5

		case tcell.KeyPgUp:
			newIndex -= 5

		case tcell.KeyEnter:
			if l.selectedFn != nil {
				l.selectedFn(l.currentIdx, l.currentItem)
			}

		case tcell.KeyEsc:
			if l.doneFn != nil {
				l.doneFn()
			}

		case tcell.KeyRune:
			switch event.Rune() {
			case ' ':
				if l.selectedFn != nil {
					l.selectedFn(l.currentIdx, l.currentItem)
				}

			case 'h':
				newIndex--

			case 'j':
				newIndex++

			case 'k':
				newIndex--

			case 'l':
				newIndex++
			}
		}

		l.SetCurrentItem(newIndex)
	})
}
