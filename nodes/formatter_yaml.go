package nodes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rivo/tview"
)

type formatterYAML struct {
	builder        strings.Builder
	indentWidth    int
	indentionCache map[int]string
	monochrome     bool
}

func newformatterYAML(intendWidth int, monochrome bool) *formatterYAML {
	return &formatterYAML{
		builder:        strings.Builder{},
		indentWidth:    intendWidth,
		indentionCache: make(map[int]string),
		monochrome:     monochrome,
	}
}

func (f formatterYAML) String() string {
	return f.builder.String()
}

func (f *formatterYAML) writeIndention(lvl int, n Node) Formatter {
	yamlLvl := lvl - 1
	if yamlLvl <= 0 {
		return f
	}

	indention, ok := f.indentionCache[yamlLvl]
	if ok == false {
		indention = strings.Repeat(" ", yamlLvl*f.indentWidth)
		f.indentionCache[yamlLvl] = indention
	}

	f.builder.WriteString(indention)

	return f
}

func (f *formatterYAML) writeDelimiter(n Node) Formatter {
	f.builder.WriteString("\n")
	return f
}

func (f *formatterYAML) writeKey(key string, n Node) Formatter {
	if f.monochrome {
		if strings.Contains(key, " ") {
			f.builder.WriteString("'")
			f.builder.WriteString(key)
			f.builder.WriteString("'")
			return f
		}
		f.builder.WriteString(key)
	} else {
		f.builder.WriteString("[lightskyblue]")
		if strings.Contains(key, " ") {
			f.builder.WriteString("'")
			f.builder.WriteString(tview.Escape(key))
			f.builder.WriteString("'")
			f.builder.WriteString("[-]")
			return f
		}
		f.builder.WriteString(tview.Escape(key))
		f.builder.WriteString("[-]")

	}
	return f
}
func (f *formatterYAML) writeKeyValueSeparator(n Node) Formatter {
	f.builder.WriteString(": ")
	return f
}

func (f *formatterYAML) writeNumber(value float64, n Node) Formatter {
	if f.monochrome {
		f.builder.WriteString(fmt.Sprintf("%g", value))
	} else {
		f.builder.WriteString(fmt.Sprintf("[darkseagreen]%g[-]", value))
	}
	return f
}
func (f *formatterYAML) writeBoolean(value bool, n Node) Formatter {
	if f.monochrome {
		f.builder.WriteString(strconv.FormatBool(value))
	} else {
		f.builder.WriteString("[deepskyblue]")
		f.builder.WriteString(strconv.FormatBool(value))
		f.builder.WriteString("[-]")
	}
	return f
}
func (f *formatterYAML) writeString(value string, n Node) Formatter {
	if f.monochrome {
		f.builder.WriteString("\"")
		f.builder.WriteString(value)
		f.builder.WriteString("\"")
	} else {
		f.builder.WriteString("[sandybrown]\"")
		f.builder.WriteString(tview.Escape(value))
		f.builder.WriteString("\"[-]")
	}
	return f
}
func (f *formatterYAML) writeNull(n Node) Formatter {
	if f.monochrome {
		f.builder.WriteString("null")
	} else {
		f.builder.WriteString("[darkskyblue]null[-]")
	}
	return f
}

func (f *formatterYAML) writeArrayItemIndicator(n Node) Formatter {
	f.builder.WriteString("- ")
	return f
}
func (f *formatterYAML) writeArrayStart(n Node) Formatter {
	if f.builder.Len() == 0 {
		return f
	}
	f.builder.WriteString("\n")
	return f
}
func (f *formatterYAML) writeArrayEnd(indentLvl int, n Node) Formatter {
	return f
}

func (f *formatterYAML) writeObjectStart(n Node) Formatter {
	if f.builder.Len() == 0 {
		return f
	}
	f.builder.WriteString("\n")
	return f
}

func (f *formatterYAML) writeObjectEnd(indentLvl int, n Node) Formatter {
	return f
}
