package widgets

import "github.com/rivo/tview"

// NewMainPage builds the promitive representing the app.
func NewMainPage(ratio int, nodeList *NodeList, output *Output, statusBar *StatusBar) (page tview.Primitive, topContentFlex *tview.Flex) {
	flex := tview.NewFlex().
		AddItem(nodeList, 0, ratio, true).
		AddItem(output, 0, 10, false)
	return tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(flex, 0, 1, true).
		AddItem(statusBar, 1, 0, false), flex
}
