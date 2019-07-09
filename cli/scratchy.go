package main

import (
	"bytes"
	"io"
	"log"
	"os"

	"github.com/jroimartin/gocui"
	"github.com/vito-c/scratchy/widgets"
)

func main() {
	jqtxt, err := os.Open("./simple.jq")
	json := os.Stdin
	var buf bytes.Buffer
	tee := io.TeeReader(json, &buf)
	logPath := "./logs/cliscratchy.log"
	if os.ExpandEnv("${CLI_SCRATCHY_LOG_FILE}") != "" {
		logPath = os.ExpandEnv("${CLI_SCRATCHY_LOG_FILE}")
	}
	f, _ := os.OpenFile(
		logPath,
		os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666,
	)
	defer f.Close()
	log.SetOutput(f)

	g, err := gocui.NewGui(gocui.Output256)

	if err != nil {
		log.Fatalln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Highlight = true
	g.SelFgColor = gocui.ColorRed
	maxX, maxY := g.Size()

	yo := int(0.3 * float64(maxY)) // yoffset
	jqbox := &widgets.JQEditor{
		Input:  "json",
		Output: "output",
		SyntaxEditor: widgets.SyntaxEditor{
			Name:     "jq",
			X:        0,
			Y:        0,
			W:        maxX - 1,
			H:        yo,
			Body:     jqtxt,
			Syntax:   "jq",
			Editable: true,
		},
	}
	jbox := &widgets.SyntaxEditor{
		Name:     "json",
		X:        0,
		Y:        yo + 1,
		W:        maxX/2 - 1,
		H:        maxY - 2 - yo,
		Body:     tee,
		Syntax:   "json",
		Editable: true,
	}
	obox := &widgets.SyntaxEditor{
		Name:     "output",
		X:        maxX / 2,
		Y:        yo + 1,
		W:        maxX / 2,
		H:        maxY - 2 - yo,
		Body:     &buf,
		Syntax:   "json",
		Editable: true,
	}
	g.SetManager(jbox, obox, jqbox)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlW, gocui.ModNone, toggleWindow); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}

var views = []string{"jq", "json", "output"}
var currentViewIndex = 0

func toggleWindow(g *gocui.Gui, v *gocui.View) error {
	nextviewIndex := 0
	if currentViewIndex < len(views)-1 {
		nextviewIndex = currentViewIndex + 1
	} else {
		nextviewIndex = 0
	}
	nextview := views[nextviewIndex]
	_, err := g.SetCurrentView(nextview)
	currentViewIndex = nextviewIndex
	return err
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
