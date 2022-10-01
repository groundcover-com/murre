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
	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	table := tview.NewTable().SetSeparator(tview.Borders.Vertical)
	app := tview.NewApplication()
	dropdown := tview.NewDropDown().
		SetLabel("Select an option (hit Enter): ").
		SetOptions([]string{"namespace", "pod", "container", "cpu", "mem"}, nil)
	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		//SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(newPrimitive("toplite"), 0, 0, 1, 1, 0, 0, false).
		AddItem(dropdown, 2, 0, 1, 1, 0, 0, false)
		// Layout for screens wider than 100 cells.
	grid.AddItem(table, 1, 0, 1, 1, 0, 100, false)
	app.SetRoot(grid, true).EnableMouse(false)
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
		cpuInMcpu := stats.CpuUsage * 1000
		return tview.NewTableCell(fmt.Sprintf("%.2fmCPU", cpuInMcpu))
	case 4:
		//convet bytes to MiB
		memoryInMiB := stats.MemoryBytes / 1024 / 1024
		return tview.NewTableCell(fmt.Sprintf("%.2fMiB", memoryInMiB))
	default:
		return nil
	}
}
