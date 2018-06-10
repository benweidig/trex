package nodes

import (
	"fmt"

	"github.com/rivo/tview"
)

type arrayNode struct {
	abstractNode
}

func (n *arrayNode) Label() Label {
	if len(n.label.text) > 0 || len(n.label.additionalInfo) > 0 {
		return n.label
	}

	if n.parent == nil {
		n.label.text = "<root>"
	} else {
		n.label.text = tview.Escape(n.key)
	}
	n.label.additionalInfo += fmt.Sprintf("[%d]", len(n.children))

	return n.label
}

func (n *arrayNode) Format(f Formatter, indentLvl int) {
	f.writeArrayStart(n)

	for idx, value := range n.children {
		if idx > 0 {
			f.writeDelimiter(n)
		}

		f.writeIndention(indentLvl, n)
		f.writeArrayItemIndicator(n)
		value.Format(f, indentLvl+1)
	}

	f.writeArrayEnd(indentLvl-1, n)
}

func (n *arrayNode) ToggleExpansion() {
	n.collapsed = !n.collapsed
}
