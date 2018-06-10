package nodes

// Label represents all info needed to print a nice node label
type Label struct {
	text           string
	additionalInfo string
}

func (l Label) String(monochrome bool) string {
	label := l.text
	if len(l.additionalInfo) > 0 {
		if len(l.text) > 0 {
			label += " "
		}
		if monochrome {
			label += l.additionalInfo
		} else {
			label += "[gray]" + l.additionalInfo + "[-]"
		}
	}
	return label
}
