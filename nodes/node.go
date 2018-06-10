package nodes

import (
	"github.com/rivo/tview"
)

// Node represents the comcept of a generic node that mostly key/value-based.
type Node interface {

	// Label contains additional info for nicer output.
	Label() Label

	// Path is the JSONPath of the node (http://goessner.net/articles/JsonPath/).
	Path() string

	// Children contains all children of the node, might be empty.
	Children() []Node

	// Format writes the formatted node (and its children) into a formatter.
	Format(f Formatter, indentLvl int)

	// IsCollapsable determinates if a node can collapse its children.
	IsCollapsable() bool

	// IsCollapsed is the current state of collapse.
	IsCollapsed() bool

	// ToggleExpansion tries to toggle the current state of expansion/collapse.
	ToggleExpansion()
}

// abstractNode is helper struct so we don't need to implement all the methods of Node
// in specialized nodes.
type abstractNode struct {
	key        string
	identifier string
	label      Label
	path       string
	parent     Node
	children   []Node
	collapsed  bool
}

// Label contains additional info for nicer output.
func (n abstractNode) Label() Label {
	return Label{
		text: tview.Escape(n.identifier),
	}
}

// Path is the JSONPath of the node (http://goessner.net/articles/JsonPath/).
func (n abstractNode) Path() string {
	return n.path
}

// Children contains all children of the node, might be empty.
func (n *abstractNode) Children() []Node {
	return n.children
}

// IsCollapsable determinates if a node can collapse its children.
func (n abstractNode) IsCollapsable() bool {
	return len(n.children) > 0
}

// IsCollapsed is the current state of collapse.
func (n abstractNode) IsCollapsed() bool {
	return n.collapsed
}

// ToggleExpansion tries to toggle the current state of expansion/collapse.
func (n *abstractNode) ToggleExpansion() {
	return
}
