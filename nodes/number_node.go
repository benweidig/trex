package nodes

type numberNode struct {
	abstractNode
	value float64
}

func (n *numberNode) Format(f Formatter, indentLvl int) {
	f.writeNumber(n.value, n)
}
