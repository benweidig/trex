package nodes

type nullNode struct {
	abstractNode
}

func (n *nullNode) Format(f Formatter, indentLvl int) {
	f.writeNull(n)
}
