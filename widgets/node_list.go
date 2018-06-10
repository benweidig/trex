package widgets

import (
	"strings"

	"github.com/benweidig/trex/nodes"
	"github.com/benweidig/trex/ui"
	"github.com/gdamore/tcell"
)

type nodeItem struct {
	node  nodes.Node
	label string
}

// Label implements interface ui.ListItem
func (n *nodeItem) Label() string {
	return n.label
}

// NodeList represents a specialized list for nodes.
type NodeList struct {
	*ui.List

	monochrome bool
	root       nodes.Node

	changedFn func(node nodes.Node)
	done      func()
}

// NewNodeList builds a new empty NodeList.
func NewNodeList(monochrome bool) *NodeList {
	n := &NodeList{
		List:       ui.NewList(),
		monochrome: monochrome,
	}
	n.SetBorder(true)

	n.setupKeyBindings()

	n.SetSelectedFn(func(index int, item ui.ListItem) {
		node := n.GetCurrentNode()
		n.toggle(node)
	})
	return n
}

// GetCurrentNode is a convience method for accessing the currently highlighted node.
func (nl *NodeList) GetCurrentNode() nodes.Node {
	item := nl.GetCurrentItem().(*nodeItem)
	return item.node
}

// SetRoot make a node the root element and generates all NodeListItems.
func (nl *NodeList) SetRoot(root nodes.Node) *NodeList {
	nl.root = root

	nl.ClearItems()

	nl.buildNodes(nl.root, 0)
	nl.SetItems(nl.GetItems(), true)

	return nl
}

func (nl *NodeList) buildNodes(node nodes.Node, indentLvl int) {
	var indention string
	if indentLvl > 0 {
		indention = strings.Repeat("  ", indentLvl)
	} else {
		indention = ""
	}
	var expansionIndicator string
	if node.IsCollapsable() {
		if node.IsCollapsed() {
			expansionIndicator = "+ "
		} else {
			expansionIndicator = "- "
		}
	} else {
		expansionIndicator = "  "
	}
	item := &nodeItem{
		label: indention + expansionIndicator + node.Label().String(nl.monochrome),
		node:  node,
	}
	nl.AddItem(item)

	if node.IsCollapsable() == false || node.IsCollapsed() || len(node.Children()) == 0 {
		return
	}

	for _, child := range node.Children() {
		nl.buildNodes(child, indentLvl+1)
	}
}

// SetChangedFn sets the handler callback on change events.
func (nl *NodeList) SetChangedFn(fn func(nodes.Node)) *NodeList {
	nl.List.SetChangedFn(func(index int, item ui.ListItem) {
		node := item.(*nodeItem)
		fn(node.node)
	})
	return nl
}

func (nl *NodeList) setupKeyBindings() {
	nl.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if nl.HasFocus() == false {
			return nil
		}

		var handled bool
		newIndex := nl.GetCurrentIdx()

		switch key := event.Key(); key {
		case tcell.KeyLeft:
			node := nl.GetCurrentNode()
			if node.IsCollapsable() && node.IsCollapsed() == false {
				nl.toggle(node)
			} else {
				newIndex--
			}
			handled = true

		case tcell.KeyRight:
			node := nl.GetCurrentNode()
			if node.IsCollapsable() && node.IsCollapsed() == true {
				nl.toggle(node)
			} else {
				newIndex++
			}
			handled = true

		case tcell.KeyRune:
			switch event.Rune() {

			case 'h':
				node := nl.GetCurrentNode()
				if node.IsCollapsable() && node.IsCollapsed() == false {
					nl.toggle(node)
				} else {
					newIndex--
				}
				handled = true

			case 'l':
				node := nl.GetCurrentNode()
				if node.IsCollapsable() && node.IsCollapsed() == true {
					nl.toggle(node)
				} else {
					newIndex++
				}
				handled = true
			}
		}
		if handled == false {
			return event
		}

		nl.SetCurrentItem(newIndex)
		return nil
	})
}

func (nl *NodeList) toggle(node nodes.Node) {
	if node.IsCollapsable() == false {
		return
	}

	node.ToggleExpansion()
	nl.SetRoot(nl.root)
}
