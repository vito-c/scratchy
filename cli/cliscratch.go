package main

import (
	"bytes"
	// "fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	// "strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/jroimartin/gocui"
	"github.com/vito-c/scratchy/scratchy/gui"
	"github.com/vito-c/scratchy/widgets"
)

type SyntaxTextBox struct {
	name     string
	x, y     int
	w, h     int
	body     io.Reader
	syntax   string
	editable bool
}

func NewSyntaxTextBox(
	name string,
	x, y int,
	w, h int,
	body io.Reader,
	syntax string,
	editable bool,
) *SyntaxTextBox {

	return &SyntaxTextBox{
		name:     name,
		x:        x,
		y:        y,
		w:        w,
		h:        h,
		body:     body,
		syntax:   syntax,
		editable: editable,
	}
}

func (b *SyntaxTextBox) Layout(g *gocui.Gui) error {

	if v, err := g.SetView(b.name, b.x, b.y, b.x+b.w, b.y+b.h); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// v.Wrap = true
		// v.Editable = true

		// if b.editable == false {
		// v.Editable = false

		style := widgets.Scratch
		if b.name == "jq" {
			log.Println("redraw jq")
		}
		if style == nil {
			log.Println("style not found falling back to default")
			style = styles.Fallback
		}

		formatter := formatters.Get("terminal256")
		if formatter == nil {
			log.Println("formatter not found falling back to default")
			formatter = formatters.Fallback
		}

		data, err := ioutil.ReadAll(b.body)
		if err != nil {
			log.Println(err.Error())
		}

		var lexer chroma.Lexer
		if b.syntax == "jq" {
			jq, _ := g.View("jq")
			ov, _ := g.View("output")
			jv, _ := g.View("json")
			ov.Editable = true

			jq.Editor = &gui.SyntaxEditor{
				Syntax: "jq",
				Input:  jv,
				Output: ov,
			}
			lexer = widgets.JqLex
		} else {
			lexer = lexers.Get(b.syntax)
			if lexer == nil {
				lexer = lexers.Fallback
			}
		}
		it, err := lexer.Tokenise(nil, string(data))
		if err != nil {
			log.Println(err.Error())
			return err
		}
		err2 := formatter.Format(v, style, it)
		if err2 != nil {
			log.Println(err.Error())
			return err
		}
		v.Editable = b.editable
		// } else {
		// 	v.Editable = true
		// }

	}

	return nil
}

func main() {
	jqtxt, err := os.Open("/Users/vito.cutten/code/personal/scratchy/simple.jq")
	json := os.Stdin
	var buf bytes.Buffer
	tee := io.TeeReader(json, &buf)
	logPath := "/Users/vito.cutten/code/personal/scratchy/logs/cliscratchy.log"
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
	jqbox := NewSyntaxTextBox("jq", 0, 0, maxX-1, yo, jqtxt, "jq", true)
	jbox := NewSyntaxTextBox("json", 0, yo+1, maxX/2-1, maxY-2-yo, tee, "json", false)
	obox := NewSyntaxTextBox("output", maxX/2, yo+1, maxX/2, maxY-2-yo, &buf, "json", false)

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

	if _, err := g.SetCurrentView("jq"); err != nil {
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
