package nodes

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rivo/tview"
)

type formatterJSON struct {
	builder        strings.Builder
	indentWidth    int
	indentionCache map[int]string
	monochrome     bool
}

func newformatterJSON(intendWidth int, monochrome bool) *formatterJSON {
	return &formatterJSON{
		builder:        strings.Builder{},
		indentWidth:    intendWidth,
		indentionCache: make(map[int]string),
		monochrome:     monochrome,
	}
}

func (f formatterJSON) String() string {
	return f.builder.String()
}

func (f *formatterJSON) writeIndention(lvl int, n Node) Formatter {
	if lvl <= 0 {
		return f
	}

	indention, ok := f.indentionCache[lvl]
	if ok == false {
		indention = strings.Repeat(" ", lvl*f.indentWidth)
		f.indentionCache[lvl] = indention
	}

	f.builder.WriteString(indention)

	return f
}

func (f *formatterJSON) writeDelimiter(n Node) Formatter {
	f.builder.WriteString(",\n")
	return f
}

func (f *formatterJSON) writeKey(key string, n Node) Formatter {
	if f.monochrome {
		f.builder.WriteString("\"")
		f.builder.WriteString(key)
		f.builder.WriteString("\"")
	} else {
		f.builder.WriteString("[lightskyblue]\"")
		f.builder.WriteString(tview.Escape(key))
		f.builder.WriteString("\"[-]")
	}
	return f
}

func (f *formatterJSON) writeKeyValueSeparator(n Node) Formatter {
	if f.monochrome {
		f.builder.WriteString(": ")
	} else {
		f.builder.WriteString("[white]:[-] ")
	}
	return f
}

func (f *formatterJSON) writeNumber(value float64, n Node) Formatter {
	if f.monochrome {
		f.builder.WriteString(fmt.Sprintf("%g", value))
	} else {
		f.builder.WriteString(fmt.Sprintf("[darkseagreen]%g[-]", value))
	}
	return f
}

func (f *formatterJSON) writeBoolean(value bool, n Node) Formatter {
	if f.monochrome {
		f.builder.WriteString(strconv.FormatBool(value))
	} else {
		f.builder.WriteString("[deepskyblue]")
		f.builder.WriteString(strconv.FormatBool(value))
		f.builder.WriteString("[-]")
	}
	return f
}

func (f *formatterJSON) writeString(value string, n Node) Formatter {
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

func (f *formatterJSON) writeNull(n Node) Formatter {
	if f.monochrome {
		f.builder.WriteString("null")
	} else {
		f.builder.WriteString("[red]null[-]")
	}
	return f
}

func (f *formatterJSON) writeArrayItemIndicator(n Node) Formatter {
	return f
}

func (f *formatterJSON) writeArrayStart(n Node) Formatter {
	f.builder.WriteString("[\n")
	return f
}

func (f *formatterJSON) writeArrayEnd(indentLvl int, n Node) Formatter {
	f.builder.WriteString("\n")
	f.writeIndention(indentLvl, n)
	f.builder.WriteString("]")
	return f
}

func (f *formatterJSON) writeObjectStart(n Node) Formatter {
	f.builder.WriteString("{\n")
	return f
}
func (f *formatterJSON) writeObjectEnd(indentLvl int, n Node) Formatter {
	f.builder.WriteString("\n")
	f.writeIndention(indentLvl, n)
	f.builder.WriteString("}")
	return f
}
