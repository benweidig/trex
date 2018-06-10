package nodes

import (
	"github.com/benweidig/trex/input"
)

// Formatter provides all functions needed to
type Formatter interface {
	String() string

	writeIndention(lvl int, n Node) Formatter
	writeDelimiter(Node) Formatter

	writeKey(string, Node) Formatter
	writeKeyValueSeparator(Node) Formatter

	writeNumber(float64, Node) Formatter
	writeBoolean(bool, Node) Formatter
	writeString(string, Node) Formatter
	writeNull(Node) Formatter

	writeArrayItemIndicator(Node) Formatter
	writeArrayStart(Node) Formatter
	writeArrayEnd(indentLvl int, n Node) Formatter

	writeObjectStart(Node) Formatter
	writeObjectEnd(indentLvl int, n Node) Formatter
}

// BuildFormatter builds the correct formatter according to its parameters.
func BuildFormatter(indentWidth int, monochrome bool, fileType input.FileType) Formatter {
	switch fileType {
	case input.FileTypeJSON:
		return newformatterJSON(indentWidth, monochrome)

	case input.FileTypeYAML:
		return newformatterYAML(indentWidth, monochrome)

	default:
		panic("No formatter for '" + string(fileType) + "'")
	}
}
