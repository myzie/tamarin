package repl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
	"github.com/mattn/go-runewidth"
	"github.com/myzie/tamarin/monkey/evaluator"
	"github.com/myzie/tamarin/monkey/lexer"
	"github.com/myzie/tamarin/monkey/object"
	"github.com/myzie/tamarin/monkey/parser"
)

const prompt = "➜ "

func createScreen(out io.Writer) tcell.Screen {
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	screen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(out, "Failed to create tcell screen: %v\n", err)
		os.Exit(1)
	}
	if err = screen.Init(); err != nil {
		fmt.Fprintf(out, "Failed to init tcell screen: %v\n", err)
		os.Exit(1)
	}
	screen.SetStyle(tcell.StyleDefault)
	screen.Clear()
	return screen
}

func emitStr(s tcell.Screen, x, y int, style tcell.Style, str string) {

	screenWidth, _ := s.Size()
	xOffset := x

	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(xOffset, y, c, comb, style)
		xOffset += w
	}

	for i := xOffset; i < screenWidth; i++ {
		s.SetContent(i, y, []rune(" ")[0], nil, style)
	}
}

func eval(env *object.Environment, line string) (string, error) {

	l := lexer.New(line)
	p := parser.New(l)
	out := new(bytes.Buffer)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		printParserErrors(out, p.Errors())
		return out.String(), errors.New("Parser error")
	}

	evaluated := evaluator.Eval(program, env)
	if evaluated != nil {
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
	return out.String(), nil
}

// Run the REPL
func Run(in io.Reader, out io.Writer) {

	white := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	red := tcell.StyleDefault.Foreground(tcell.ColorRed)
	green := tcell.StyleDefault.Foreground(tcell.ColorLawnGreen)

	screen := createScreen(out)
	screen.EnableMouse()
	screen.SetStyle(white)

	env := object.NewEnvironment()

	var err error
	var input, output string

	running := true
	for running {

		emitStr(screen, 1, 1, green, "Tamarin Cloud Shell")

		row := 2

		// Show current input
		for _, line := range strings.Split(input, "\n") {
			emitStr(screen, 1, row, white, line)
			row++
		}

		// Show cursor
		screen.ShowCursor(len(input)+1, row-1)

		// Show current error
		if err != nil {
			for _, line := range strings.Split(err.Error(), "\n") {
				emitStr(screen, 1, row, red, line)
				row++
			}
		}

		// Show current output
		for _, line := range strings.Split(output, "\n") {
			emitStr(screen, 1, row, white, line)
			row++
		}

		// Show environment
		envRow := 1
		screenWidth, _ := screen.Size()
		for _, key := range env.Keys() {
			value, _ := env.Get(key)
			displayStr := fmt.Sprintf("%s: %v", key, value.Inspect())
			emitStr(screen, screenWidth/2, envRow, white, displayStr)
			envRow++
		}

		screen.Show()

		// Wait for the next input event
		event := screen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyEscape, tcell.KeyCtrlC:
				running = false
				break
			case tcell.KeyEnter:
				output, err = eval(env, input)
				input = ""
			case tcell.KeyBackspace, tcell.KeyBackspace2, tcell.KeyDelete:
				lineLen := len(input)
				if lineLen > 0 {
					input = input[0 : lineLen-1]
				}
			default:
				if strconv.IsPrint(event.Rune()) {
					// Clear error and output when user starts entering a new
					// command
					if input == "" {
						err = nil
						output = ""
						screen.Clear()
					}
					input += string(event.Rune())
				}
			}
		case *tcell.EventResize:
			screen.Sync()
		}
	}
	screen.Fini()
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "Woops!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
