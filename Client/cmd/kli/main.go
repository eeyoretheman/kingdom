package main

import (
	readwriter "kli/internal/readwriter"
	"log"
	"net"
	"os"
	"strings"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	tb "github.com/nsf/termbox-go"
)

const (
	UpdateTellers int = 1
	UpdateHistory int = 2
)

type Data struct {
	history      map[string][]string
	server       string
	tellers      []string
	input        string
	cursor       int
	selection    int
	teller       string
	historyShift int
}

type State struct {
	history *widgets.List
	server  *widgets.Paragraph
	teller  *widgets.List
	input   *widgets.Paragraph
}

func initialize() State {
	tb.SetInputMode(tb.InputEsc)
	width, height := ui.TerminalDimensions()

	server := widgets.NewParagraph()
	server.SetRect(0, 0, width, 3)
	server.Title = "Server"

	teller := widgets.NewList()
	teller.WrapText = true
	teller.Title = "Target"

	teller.SetRect(0, 3, width, 6)

	history := widgets.NewList()
	history.WrapText = true

	history.SetRect(0, 6, width, height-3)

	input := widgets.NewParagraph()
	input.SetRect(0, height-3, width, height)

	return State{history: history, server: server, teller: teller, input: input}
}

func render(state State, data Data) {
	width, height := ui.TerminalDimensions()

	state.server.SetRect(0, 0, width, 3)
	state.teller.SetRect(0, 3, width, 6)
	state.history.SetRect(0, 6, width, height-3)
	state.input.SetRect(0, height-3, width, height)

	state.server.Text = data.server
	state.teller.Rows = data.tellers

	state.history.Rows = data.history[data.teller][data.historyShift:]
	state.input.Text = data.input

	switch data.selection {
	case 1:
		state.teller.SetRect(0, 3, width, 9)
		state.history.SetRect(0, 9, width, height-3)
		state.teller.BorderStyle.Fg = ui.ColorGreen
		state.server.BorderStyle.Fg = ui.ColorWhite
		state.input.BorderStyle.Fg = ui.ColorWhite
		state.teller.SelectedRowStyle.Bg = ui.ColorGreen
	case 2:
		state.input.Text = state.input.Text[:data.cursor] + "_" + state.input.Text[data.cursor:]
		state.input.BorderStyle.Fg = ui.ColorGreen
		state.server.BorderStyle.Fg = ui.ColorWhite
		state.teller.BorderStyle.Fg = ui.ColorWhite
		state.teller.SelectedRowStyle.Bg = ui.ColorClear
	}

	ui.Render(state.server, state.teller, state.history, state.input)
}

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("Could not initialize termui: %v", err)
	}
	defer ui.Close()

	var server string

	if len(os.Args) > 1 {
		server = os.Args[1]
	} else {
		server = "127.0.0.1:2222"
	}

	socket, err := net.Dial("tcp", server)

	if err != nil {
		ui.Close()
		log.Fatalf("Could not connect to server: %v", err)
	}

	requests := make(chan readwriter.Request)
	responses := make(chan []byte)

	connection := readwriter.ReadWriter{Conn: socket, From: "!", Request: requests, Response: responses}
	go readwriter.ReadWriterHandler(&connection)

	id := strings.TrimSuffix(string(<-connection.Response), "\n")
	connection.From = id

	state := initialize()

	history := make(map[string][]string)
	//history["!"] = []string{}

	data := Data{history: history, tellers: []string{"!"}, server: server, cursor: 0, selection: 2, teller: "!", historyShift: 0}

	commandHistory := make([]string, 0)
	index := 0

	//data.tellers = updateTellers(&connection)

	connection.Request <- readwriter.Request{To: "!", Command: "lst", Body: "."}

	render(state, data)

	uiEvents := ui.PollEvents()

	for {
		select {
		case r := <-connection.Response:
			text := string(r)

			if strings.Contains(text, "#!!#") {
				text = strings.ReplaceAll(text, "#!!#", "")
				data.tellers = []string{"!"}

				if text == "No tellers\n" {
					break
				}

				data.tellers = []string{"!"}
				data.tellers = append(data.tellers, strings.Split(text, "\n")...)
			} else {
				data.history[data.teller] = append(data.history[data.teller], text)
			}
		case e := <-uiEvents:
			switch data.selection {
			case 2:
				switch e.ID {
				case "<C-c>":
					return
				case "<Backspace>":
					if len(data.input) > 0 {
						data.input = data.input[:data.cursor-1] + data.input[data.cursor:]
						data.cursor -= 1
					}
				case "<Space>":
					if data.cursor == len(data.input) {
						data.input += " "
					} else {
						data.input = data.input[:data.cursor] + " " + data.input[data.cursor:]
					}

					data.cursor += 1
				case "<Resize>", "<MouseLeft>", "<MouseRight>", "<MouseRelease>", "<Tab>":
				case "<Left>":
					if data.cursor > 0 {
						data.cursor -= 1
					}
				case "<Right>":
					if data.cursor < len(data.input) {
						data.cursor += 1
					}
				case "<Up>":
					if index > 0 {
						index -= 1
						data.input = commandHistory[index]
						data.cursor = len(commandHistory[index])
					}
				case "<Down>":
					if index < len(commandHistory)-1 {
						index += 1
						data.input = commandHistory[index]
						data.cursor = len(commandHistory[index])
					}
				case "<Enter>":
					data.history[data.teller] = append(data.history[data.teller], "# "+data.input)
					commandHistory = append(commandHistory, data.input)
					index += 1

					parts := strings.SplitN(data.input, " ", 2)

					switch len(parts) {
					case 1:
						connection.Request <- readwriter.Request{To: data.teller, Command: parts[0], Body: ".\n"}
					case 2:
						connection.Request <- readwriter.Request{To: data.teller, Command: parts[0], Body: parts[1] + "\n"}
					}

					data.input = ""
					data.cursor = 0

					connection.Request <- readwriter.Request{To: "!", Command: "lst", Body: "."}
				case "<C-u>":
					if data.selection > 1 {
						data.selection -= 1
					}

					connection.Request <- readwriter.Request{To: "!", Command: "lst", Body: "."}
				case "<C-y>":
					if data.selection < 2 {
						data.selection += 1
					}

					connection.Request <- readwriter.Request{To: "!", Command: "lst", Body: "."}
				case "<C-j>":
					if data.historyShift < len(data.history[data.teller])-1 {
						data.historyShift += 1
					}
				case "<C-k>":
					if data.historyShift > 0 {
						data.historyShift -= 1
					}
				default:
					if data.cursor == len(data.input) {
						data.input += e.ID
					} else {
						data.input = data.input[:data.cursor] + e.ID + data.input[data.cursor:]
					}

					data.cursor += 1
				}
			case 1:
				switch e.ID {
				case "<C-c>":
					return
				case "<C-u>":
					if data.selection > 1 {
						data.selection -= 1
					}
				case "<C-y>":
					if data.selection < 2 {
						data.selection += 1
					}
				case "<Down>":
					state.teller.ScrollDown()
				case "<Up>":
					state.teller.ScrollUp()
				case "<Enter>":
					data.teller = strings.Split(state.teller.Rows[state.teller.SelectedRow], ",")[0]

					if len(data.history[data.teller]) < 5 {
						data.historyShift = 0
					} else {
						data.historyShift = len(data.history[data.teller]) - 5
					}

					data.selection = 2
				}
			}
		}

		ui.Clear()
		render(state, data)
	}
}
