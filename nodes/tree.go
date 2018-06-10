package nodes

import (
	"fmt"
	"sort"
	"strings"

	"github.com/benweidig/trex/input"
)

// Tree represents a tree of nodes with a single root
type Tree struct {
	fileType input.FileType
	root     Node
}

// NewTree builds a new tree
func NewTree(fileType input.FileType, raw interface{}) (*Tree, error) {
	t := &Tree{
		fileType: fileType,
	}

	root, err := buildNodes("$", "", "", nil, raw)
	if err != nil {
		return nil, err
	}
	t.root = root

	return t, nil
}

// FileType return the input.FileType of the tree
func (t Tree) FileType() input.FileType {
	return t.fileType
}

// Root returns the root node of the tree
func (t *Tree) Root() Node {
	return t.root
}

func buildNodes(path string, key string, identifier string, parent Node, j interface{}) (Node, error) {
	var node Node

	switch value := j.(type) {
	case bool:
		node = &boolNode{
			abstractNode{
				key:        key,
				identifier: identifier,
				path:       path,
				parent:     parent,
			},
			value,
		}

	case string:
		node = &stringNode{
			abstractNode{
				key:        key,
				identifier: identifier,
				path:       path,
				parent:     parent,
			},
			safeString(value),
		}

	case float64:
		node = &numberNode{
			abstractNode{
				key:        key,
				identifier: identifier,
				path:       path,
				parent:     parent,
			},
			value,
		}

	case map[interface{}]interface{}:
		objectNode := &objectNode{
			abstractNode: abstractNode{
				key:        key,
				identifier: identifier,
				path:       path,
				parent:     parent,
				children:   make([]Node, len(value)),
			},
			values: make(map[string]Node),
		}

		keys := make([]string, len(value))

		idx := 0
		for key := range value {
			keys[idx] = key.(string)
			idx++
		}

		sort.Strings(keys)

		for idx, childKey := range keys {
			childInterface := value[childKey]
			childPath := path + "." + childKey
			childNode, err := buildNodes(childPath, childKey, childKey, objectNode, childInterface)
			if err != nil {
				return nil, err
			}
			objectNode.values[childKey] = childNode
			objectNode.children[idx] = childNode
		}
		node = objectNode

	case map[string]interface{}:
		objectNode := &objectNode{
			abstractNode: abstractNode{
				key:        key,
				identifier: identifier,
				path:       path,
				parent:     parent,
				children:   make([]Node, len(value)),
			},
			values: make(map[string]Node),
		}

		keys := make([]string, len(value))

		idx := 0
		for key := range value {
			keys[idx] = key
			idx++
		}

		sort.Strings(keys)

		for idx, childKey := range keys {
			childInterface := value[childKey]
			childPath := path + "." + childKey
			childNode, err := buildNodes(childPath, childKey, childKey, objectNode, childInterface)
			if err != nil {
				return nil, err
			}
			objectNode.values[childKey] = childNode
			objectNode.children[idx] = childNode
		}
		node = objectNode

	case []interface{}:
		arrayNode := &arrayNode{
			abstractNode{
				key:        key,
				identifier: identifier,
				path:       path + "[*]",
				parent:     parent,
				children:   make([]Node, len(value)),
			},
		}
		for idx, childInterface := range value {
			arrayIndex := fmt.Sprintf("[%d]", idx)
			childNode, err := buildNodes(path+arrayIndex, "", arrayIndex, arrayNode, childInterface)
			if err != nil {
				return nil, err
			}
			// Values and Children is identical, optimize in the future
			arrayNode.children[idx] = childNode
		}
		node = arrayNode
	default:
		if value != nil {
			panic("Invalid type, path = " + path)
		}
		node = &nullNode{
			abstractNode{
				key:        key,
				identifier: identifier,
				parent:     parent,
				path:       path,
			},
		}
	}
	return node, nil
}

func safeString(str string) string {
	safe := str
	safe = strings.Replace(safe, "\n", "\\n", -1)
	safe = strings.Replace(safe, "\r", "\\r", -1)
	safe = strings.Replace(safe, "\t", "\\t", -1)
	return safe
}
