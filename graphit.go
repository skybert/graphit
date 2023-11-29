// Program that reads a PID as the first argument and displays a
// real-time TUI graph of the memory usage of that process.
//
// Author: torstein at skybert.net

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/struCoder/pidusage"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func main() {
	if err := ui.Init(); err != nil {
		log.Fatalf("Failed initializing termui: %v", err)
	}
	defer ui.Close()

	if len(os.Args) < 2 {
		fmt.Println("No process ID provided")
		os.Exit(1)
	}
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid process ID provided")
		os.Exit(1)
	}

	data := []float64{}

	sl0 := widgets.NewSparkline()
	sl0.Data = data
	sl0.LineColor = ui.ColorGreen

	// single
	slg0 := widgets.NewSparklineGroup(sl0)
	slg0.Title = "Memory usage of PID " + strconv.Itoa(pid)

	width, height := ui.TerminalDimensions()
	slg0.SetRect(0, 0, width, height)

	// text box with info about the command
	p := widgets.NewParagraph()
	pPadding := 2
	p.Title = "Command"
	p.Text = getCommandName(pid)
	p.SetRect(pPadding, (height - 5), (width - pPadding), (height - 10))
	p.TextStyle.Fg = ui.ColorWhite
	p.BorderStyle.Fg = ui.ColorCyan

	grid := ui.NewGrid()
	grid.SetRect(0, 0, width, height)
	grid.Set(
		ui.NewRow(2.0/3,
			ui.NewCol(1.0, slg0)),
		ui.NewRow(1.0/3, ui.NewCol(1.0, p)))

	ticker := time.NewTicker(time.Second).C

	draw := func() {
		memoryUsage := getMemoryProcessOfProcess(pid)
		if memoryUsage > 0 {
			slg0.Title = getTitle(pid, memoryUsage)
			sl0.Data = suffleToTheLeft(sl0.Data, memoryUsage)
		} else {
			slg0.Title = "PID " + strconv.Itoa(pid) + " terminated"
		}

		ui.Render(grid)
	}

	uiEvents := ui.PollEvents()

	for {
		select {
		case e := <-uiEvents:
			switch e.Type {
			case ui.KeyboardEvent:
				if e.ID == "q" || e.ID == "<C-c>" || e.ID == "<Escape>" {
					// ensure terminal isn't
					// messed up when the
					// application exits
					defer ui.Clear()
					return
				}
			case ui.ResizeEvent:
				payload := e.Payload.(ui.Resize)
				grid.SetRect(0, 0, payload.Width, payload.Height)
				ui.Clear()
				ui.Render(grid)
			}
		case <-ticker:
			// Update the gauge value with the current memory usage
			draw()
		}
	}
}

func getCommandName(pid int) string {
	cmd, err := os.ReadFile("/proc/" + strconv.Itoa(pid) + "/cmdline")
	if err != nil {
		return "unknown"
	}

	return strings.Join(strings.Split(string(cmd), "\x00"), " ")
}

func getTitle(pid int, memoryUsage float64) string {
	return "PID " + strconv.Itoa(pid) + " uses " + strconv.FormatFloat(memoryUsage, 'f', 0, 64) + " bytes"
}

func suffleToTheLeft(data []float64, newValue float64) []float64 {
	width, _ := ui.TerminalDimensions()

	numberThatCanFitWithinTerminalWidth := width - 1
	data = append(data, newValue)
	if len(data) == 1 {
		return data
	} else if len(data) < numberThatCanFitWithinTerminalWidth {
		return data
	} else {
		return data[1:]
	}
}

func getMemoryProcessOfProcess(pid int) float64 {
	sysInfo, err := pidusage.GetStat(pid)
	if err != nil {
		return -1
	}

	return sysInfo.Memory
}
