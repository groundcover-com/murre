package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Table struct {
	app   *tview.Application
	table *tview.Table
}

func CreateNewTable() *Table {
	table := tview.NewTable().SetSeparator(tview.Borders.Vertical)
	app := tview.NewApplication()
	app.SetRoot(table, true).EnableMouse(false)
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEscape ||
			event.Key() == tcell.KeyCtrlC ||
			event.Rune() == 'Q' ||
			event.Rune() == 'q' {
			app.Stop()
		}
		return event
	})
	return &Table{
		app:   app,
		table: table,
	}
}

func (t *Table) Draw() error {
	return t.app.Run()
}

func (t *Table) Update(stats []*Stats) {
	t.app.QueueUpdateDraw(func() {
		t.table.Clear()
		t.updateColumns()
		for i, stat := range stats {
			for j := 0; j < 5; j++ {
				t.table.SetCell(i+1, j, t.getCell(stat, j).SetExpansion(1))
			}
		}
		t.table.ScrollToBeginning()

	})
}

func (t *Table) updateColumns() {
	yel := tcell.ColorAliceBlue
	t.table.SetCell(0, 0, tview.NewTableCell("namespace").SetExpansion(1).SetAlign(tview.AlignCenter)).SetTitleColor(yel)
	t.table.SetCell(0, 1, tview.NewTableCell("pod").SetExpansion(1).SetAlign(tview.AlignCenter)).SetTitleColor(yel)
	t.table.SetCell(0, 2, tview.NewTableCell("container").SetExpansion(1).SetAlign(tview.AlignCenter)).SetTitleColor(yel)
	t.table.SetCell(0, 3, tview.NewTableCell("cpu").SetExpansion(1).SetAlign(tview.AlignCenter)).SetTitleColor(yel)
	t.table.SetCell(0, 4, tview.NewTableCell("mem").SetExpansion(1).SetAlign(tview.AlignCenter)).SetTitleColor(yel)
}

func (t *Table) getCell(stats *Stats, column int) *tview.TableCell {
	switch column {
	case 0:
		return tview.NewTableCell(stats.Namespace)
	case 1:
		return tview.NewTableCell(stats.PodName)
	case 2:
		return tview.NewTableCell(stats.ContainerName)
	case 3:
		if stats.CpuUsage <= 0 {
			return tview.NewTableCell("-").SetAlign(tview.AlignCenter)
		}
		cpuInMcpu := stats.CpuUsage * 1000
		return tview.NewTableCell(fmt.Sprintf("%.2fmCPU", cpuInMcpu))
	case 4:
		if stats.MemoryBytes <= 0 {
			return tview.NewTableCell("-").SetAlign(tview.AlignCenter)
		}
		//convet bytes to MiB
		memoryInMiB := stats.MemoryBytes / 1024 / 1024
		return tview.NewTableCell(fmt.Sprintf("%.2fMiB", memoryInMiB))
	default:
		return nil
	}
}
