package widgets

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/vito-c/scratchy/jq"
)

var jqr *jq.JQ

type JQEditor struct {
	SyntaxEditor
	Input  string
	Output string
}

func run(query string) {
	jq.Init()
	jq.Path = "/usr/local/bin/jq"
	jqr = &jq.JQ{
		J: input.Buffer(),
		Q: query,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if jqr.ValidateFilter() == nil {
		log.Println("filter is valid")
		output.Clear()
		i := input.Buffer()
		log.Println("input: ", len(i))
		if err := jqr.Eval(ctx, output, ioutil.Discard); err != nil {
			log.Println("err: ", err.Error())
		}
	}
}

var output *gocui.View
var input *gocui.View
var gui *gocui.Gui

func (e *JQEditor) Layout(g *gocui.Gui) error {
	if v, err := g.SetView(e.Name, e.X, e.Y, e.X+e.W, e.Y+e.H); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = e.Editable
		v.Editor = e

		if data, err := ioutil.ReadAll(e.Body); err == nil {
			log.Println(e.Name, ": update layout")
			updateColors(e.Syntax, v, string(data))
		}
		output, _ = g.View(e.Output)
		input, _ = g.View(e.Input)
		gui = g
		
		data := v.Buffer()
		run(data)
		updateColors(e.Syntax, v, data)
		updateColors("json", output, output.Buffer())

	}

	return nil
}

func (ve *JQEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
		data := v.Buffer()
		if output == nil {
			output.Rewind()
			output, _ = gui.View(ve.Output)
		}
		if input == nil {
			input.Rewind()
			input, _ = gui.View(ve.Input)
		}
		run(data)
		updateColors(ve.Syntax, v, data)
		updateColors("json", output, output.Buffer())
		// log.Println(v.Name(), ": update colors")
	case key == gocui.KeySpace:
		v.EditWrite(' ')
		// ve.Run()
		// ve.updateColors(v, v.Buffer())
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyInsert:
		v.Overwrite = !v.Overwrite
		log.Println("overwrite")
	case key == gocui.KeyEnter:
		v.EditNewLine()
	case key == gocui.KeyCtrlD:
		v.MoveCursor(0, 1, false)
	case key == gocui.KeyCtrlU:
		v.MoveCursor(0, -1, false)
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
