package nodes

type boolNode struct {
	abstractNode
	value bool
}

func (n *boolNode) Format(f Formatter, indentLvl int) {
	f.writeBoolean(n.value, n)
}
