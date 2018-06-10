package nodes

type stringNode struct {
	abstractNode
	value string
}

func (n *stringNode) Format(f Formatter, indentLvl int) {
	f.writeString(n.value, n)
}
