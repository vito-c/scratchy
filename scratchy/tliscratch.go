package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/alecthomas/chroma/formatters"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/rivo/tview"
)

var json *os.File

func main() {
	json = os.Stdin

	logPath := "/Users/vito.cutten/code/personal/scratchy/logs/tliscratchy.log"
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

	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(false).
		SetRegions(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	go func() {
		setRichText(textView, json)
	}()

	textView.SetBorder(true)
	log.Println("in main loop")
	if err := app.SetRoot(textView, true).SetFocus(textView).Run(); err != nil {
		panic(err)
	}
}

func setRichText(w io.Writer, r io.Reader) error {
	style := styles.Get("monokai")
	if style == nil {
		log.Println("style not found falling back to default")
		style = styles.Fallback
	}
	formatter := formatters.Get("terminal16m")
	if formatter == nil {
		log.Println("formatter not found falling back to default")
		formatter = formatters.Fallback
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Println(err.Error())
	}
	lexer := lexers.Get("json")
	if lexer == nil {
		lexer = lexers.Fallback
	}
	log.Println("bytes len: ", len(b))
	it, err := lexer.Tokenise(nil, string(b))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	log.Println("setting text")
	err2 := formatter.Format(w, style, it)
	log.Println("text set")
	if err2 != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
