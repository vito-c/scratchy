package gui

import (
	"context"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/jroimartin/gocui"
	"github.com/vito-c/scratchy/jq"
	"github.com/vito-c/scratchy/widgets"
)

var jqr *jq.JQ

type SyntaxEditor struct {
	Syntax string
	Input  *gocui.View
	Output *gocui.View
	// G gocui.Gui
}

func (ve *SyntaxEditor) updateColors(qv *gocui.View) {
	jq.Init()
	jq.Path = "/usr/local/bin/jq"
	style := widgets.Scratch
	if style == nil {
		log.Println("style not found falling back to default")
		style = styles.Fallback
	}

	formatter := formatters.Get("terminal256")
	if formatter == nil {
		log.Println("formatter not found falling back to default")
		formatter = formatters.Fallback
	}
	var lexer chroma.Lexer
	if ve.Syntax == "jq" {
		lexer = widgets.JqLex
	} else {
		lexer = lexers.Get(ve.Syntax)
		if lexer == nil {
			log.Println("Using fallback lexer")
			lexer = lexers.Fallback
		}
	}

	// .Update(func(g *gocui.Gui) error {
	log.Println("Update")
	// 	ov, _ := g.View("output")
	// 	jv, _ := g.View("json")
	// 	ov.Buffer()
	jqr = &jq.JQ{
		J: "",
		Q: string(qv.Buffer()),
	}
	// if err := jqr.Validate(); err != nil {
	// 	log.Println("ERR: ", err.Error())
	// 	// return nil
	// }
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// ov.Clear()

	if jqr.ValidateFilter() == nil {
		log.Println("filter is valid")
		ve.Output.Clear()
		ve.Input.Buffer()
		if err := jqr.EvalStream(ctx, strings.NewReader(ve.Input.Buffer()), ve.Output, ioutil.Discard); err != nil {
			log.Println(err.Error())
		}
	}
	it, _ := lexer.Tokenise(nil, qv.Buffer())
	qv.Clear()
	err := formatter.Format(qv, style, it)
	if err != nil {
		log.Println(err.Error())
	}

	// 	return nil
	// })
}

func (ve *SyntaxEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	log.Println(key)
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
		ve.updateColors(v)
		log.Println("update colors")
	case key == gocui.KeySpace:
		v.EditWrite(' ')
		ve.updateColors(v)
		log.Println("update colors")
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
		log.Println("overwrite")
	case key == gocui.KeyEnter:
		v.EditNewLine()
	case key == gocui.KeyArrowDown:
		v.MoveCursor(0, 1, false)
	case key == gocui.KeyArrowUp:
		v.MoveCursor(0, -1, false)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
}
