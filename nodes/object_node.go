package nodes

import (
	"fmt"
	"sort"

	"github.com/rivo/tview"
)

type objectNode struct {
	abstractNode

	values map[string]Node
}

func (n *objectNode) Label() Label {
	if len(n.label.text) > 0 || len(n.label.additionalInfo) > 0 {
		return n.label
	}
	if n.parent == nil {
		n.label.text = "<root>"
	} else {
		n.label.text = tview.Escape(n.key)
	}

	if len(n.label.text) == 0 {
		n.label.text = n.identifier
	}
	n.label.additionalInfo = fmt.Sprintf("{%d}", len(n.values))

	return n.label
}

func (n *objectNode) Format(f Formatter, indentLvl int) {
	keys := make([]string, len(n.values))

	idx := 0
	for key := range n.values {
		keys[idx] = key
		idx++
	}

	sort.Strings(keys)

	f.writeObjectStart(n)

	for idx, key := range keys {
		if idx > 0 {
			f.writeDelimiter(n)
		}

		value, _ := n.values[key]
		f.writeIndention(indentLvl, n)
		f.writeKey(key, n)
		f.writeKeyValueSeparator(n)

		value.Format(f, indentLvl+1)
	}

	f.writeObjectEnd(indentLvl-1, n)
}

func (n *objectNode) ToggleExpansion() {
	n.collapsed = !n.collapsed
}
