package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/benweidig/trex/input"
	"github.com/benweidig/trex/nodes"
	"github.com/benweidig/trex/ui"
	"github.com/benweidig/trex/version"
	"github.com/benweidig/trex/widgets"
	"github.com/rivo/tview"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

var (
	monochromeArg bool
)

// RootCmd is the only command, so this is t-rex
var RootCmd = &cobra.Command{
	Version: version.BuildVersion(),
	Use:     "trex",
	Short:   "Tree Explorer",
	Args:    cobra.MaximumNArgs(1),
	Long:    "CLI tool for visualizing JSON/YAML files",
	Run:     runCommand,
}

const (
	uiHelpPage = "ui.page.help"
)

var (
	app               = tview.NewApplication()
	pages             = tview.NewPages()
	uiNodeList        *widgets.NodeList
	uiOutput          *widgets.Output
	uiStatusBar       = widgets.NewStatusBar()
	leftToRightRatio  = 3
	formatterFileType input.FileType
	topContentFlex    *tview.Flex
)

func init() {
	ui.ApplyStyling()
	RootCmd.Flags().BoolVarP(&monochromeArg, "monochrome", "m", false, "Monochrome output, no ANSI colors")
}

func runCommand(_ *cobra.Command, args []string) {
	bytes, fileType, err := getBytes(args)
	if err != nil {
		panic(err)
	}

	raw, err := input.Load(fileType, bytes)
	tree, err := nodes.NewTree(fileType, raw)
	if err != nil {
		panic(err)
	}

	// Check if the terminal actually supports colors
	monochromeArg = monochromeArg || os.Getenv("TERM") == "dumb" ||
		(!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))
	_, noColorExists := os.LookupEnv("NO_COLOR")
	if monochromeArg == false && noColorExists {
		monochromeArg = true
	}

	uiOutput = widgets.NewOutput(monochromeArg)

	uiNodeList = widgets.NewNodeList(monochromeArg)
	uiNodeList.SetChangedFn(func(node nodes.Node) {
		if formatterFileType == input.FileTypeUnknown {
			formatterFileType = tree.FileType()
		}
		f := nodes.BuildFormatter(2, monochromeArg, formatterFileType)
		node.Format(f, 1)
		uiOutput.SetText(f.String())
		uiStatusBar.SetContent(node.Path(), formatterFileType)
	})
	uiNodeList.SetRoot(tree.Root())

	outputPopup := widgets.NewFormatterPopup(func(selected input.FileType) {
		formatterFileType = selected
		uiNodeList.TriggerChanged()
		pages.SwitchToPage(widgets.MainPage)
		app.SetFocus(uiNodeList)
	})

	helpPopup := widgets.NewHelpPopup()

	mainPage, topContentFlex := widgets.NewMainPage(leftToRightRatio, uiNodeList, uiOutput, uiStatusBar)

	pages.
		AddPage(widgets.MainPage, mainPage, true, true).
		AddPage(widgets.FormatterPopupPage, outputPopup, true, false).
		AddPage(widgets.HelpPopupPage, helpPopup, true, false)

	app.
		SetRoot(pages, true).
		SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

			// Most events should only be accessible when no popup is active
			if mainPage.GetFocusable().HasFocus() {

				switch event.Key() {
				case tcell.KeyRune:
					switch event.Rune() {
					case 'f': // Choose Formatter
						pages.ShowPage(widgets.FormatterPopupPage)
						app.SetFocus(outputPopup)
						app.Draw()
						return nil

					case '?': // Help
						pages.ShowPage(widgets.HelpPopupPage)
						app.SetFocus(helpPopup)
						app.Draw()
						return nil

					case 'c': // Copy
						f := nodes.BuildFormatter(2, true, formatterFileType)
						node := uiNodeList.GetCurrentNode()
						node.Format(f, 1)
						err := clipboard.WriteAll(f.String())
						if err != nil {
							panic(err)
						}
						app.Draw()
						return nil
					}

				case tcell.KeyTab:
					if uiNodeList.HasFocus() {
						app.SetFocus(uiOutput)
					} else {
						app.SetFocus(uiNodeList)
					}
					app.Draw()

					return nil

				case tcell.KeyLeft:
					if mainPage.GetFocusable().HasFocus() == false {
						return event
					}
					if event.Modifiers() != tcell.ModShift || outputPopup.GetFocusable().HasFocus() {
						break
					}
					if leftToRightRatio < 2 {
						return nil
					}
					leftToRightRatio--

					topContentFlex.ResizeItem(uiNodeList, 0, leftToRightRatio)
					app.Draw()
					return nil

				case tcell.KeyRight:
					if event.Modifiers() != tcell.ModShift || outputPopup.GetFocusable().HasFocus() {
						break
					}
					if leftToRightRatio > 6 {
						return nil
					}
					leftToRightRatio++
					topContentFlex.ResizeItem(uiNodeList, 0, leftToRightRatio)
					app.Draw()
					return nil
				}

				return event
			}

			switch event.Key() {
			case tcell.KeyEsc:
				pages.SwitchToPage(widgets.MainPage)
				outputPopup.Blur()
				app.SetFocus(uiNodeList)
				app.Draw()
				return nil
			}

			return event
		})

	app.Run()
}

const askIfBiggerThanMB = 20

func getBytes(args []string) ([]byte, input.FileType, error) {
	var fileType input.FileType
	// Piped in content wins over file
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		// Read from stdin
		bytes, err := ioutil.ReadAll(os.Stdin)
		return bytes, fileType, err
	}

	if len(args) != 1 {
		return nil, fileType, errors.New("No file specified and now piped input detected")
	}
	path := args[0]

	fi, err := os.Stat(path)
	if err != nil {
		return nil, input.FileTypeUnknown, err
	}
	sizeMB := fi.Size() / 1024 / 1024
	if sizeMB > askIfBiggerThanMB {
		proceed, err := askQuestionYN(fmt.Sprintf("JSON file > %d MB! Trex might eat up all CPU/RAM. Proceed?", askIfBiggerThanMB))
		if err != nil {
			return nil, input.FileTypeUnknown, err
		}
		if proceed == false {
			os.Exit(0)
		}
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return bytes, fileType, errors.New("Could load file")
	}

	fileType = input.DetectFileType(path)
	return bytes, fileType, err
}
