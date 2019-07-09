package widgets

import (
	"io"
	"io/ioutil"
	"log"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/jroimartin/gocui"
)

type SyntaxEditor struct {
	Syntax   string
	Name     string
	X, Y     int
	W, H     int
	Body     io.Reader
	Editable bool
}

func (e *SyntaxEditor) Layout(g *gocui.Gui) error {

	if v, err := g.SetView(e.Name, e.X, e.Y, e.X+e.W, e.Y+e.H); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = e.Editable
		v.Editor = e

		if data, err := ioutil.ReadAll(e.Body); err == nil {
			log.Println(e.Name, ": layout")
			updateColors(e.Syntax, v, string(data))
		}
	}

	return nil
}

func updateColors(syntax string, v *gocui.View, data string) {
	style := Scratch
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
	if syntax == "jq" {
		lexer = JqLex
	} else {
		lexer = lexers.Get(syntax)
		if lexer == nil {
			log.Println("Using fallback lexer")
			lexer = lexers.Fallback
		}
	}

	log.Println(v.Name(), ": data: ", len(data))
	if it, err := lexer.Tokenise(nil, data); err == nil {
		v.Clear()
		err := formatter.Format(v, style, it)
		if err != nil {
			log.Println(err.Error())
		}
	}

}

func (ve *SyntaxEditor) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	log.Println(key)
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
		updateColors(ve.Syntax, v, v.Buffer())
		log.Println(v.Name(), ": update colors")
	case key == gocui.KeySpace:
		v.EditWrite(' ')
		updateColors(ve.Syntax, v, v.Buffer())
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
