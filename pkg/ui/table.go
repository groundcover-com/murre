package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/groundcover-com/murre/pkg/k8s"

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

func (t *Table) Update(stats []*k8s.Stats) {
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
	blue := tcell.ColorBlue
	t.table.SetCell(0, 0, t.createColumnCell("Namespace").SetTextColor(blue))
	t.table.SetCell(0, 1, t.createColumnCell("Pod").SetTextColor(blue))
	t.table.SetCell(0, 2, t.createColumnCell("Container").SetTextColor(blue))
	t.table.SetCell(0, 3, t.createColumnCell("CPU").SetTextColor(blue))
	t.table.SetCell(0, 4, t.createColumnCell("Memory").SetTextColor(blue))
}

func (t *Table) createColumnCell(text string) *tview.TableCell {
	return tview.NewTableCell(text).SetAlign(tview.AlignCenter).SetTextColor(tcell.ColorBlue).SetBackgroundColor(tcell.ColorDarkGray)
}

func (t *Table) getCell(stats *k8s.Stats, column int) *tview.TableCell {
	switch column {
	case 0:
		return tview.NewTableCell(stats.Namespace)
	case 1:
		return tview.NewTableCell(stats.PodName)
	case 2:
		return tview.NewTableCell(stats.ContainerName)
	case 3:
		if stats.CpuUsageMilli <= 0 {
			return tview.NewTableCell("\u23F1").SetAlign(tview.AlignCenter)
		}
		if stats.CpuUsagePercent > 0 {
			color := t.getCellColor(stats.CpuUsagePercent)
			return tview.NewTableCell(fmt.Sprintf("%.0f/%.0fmCPU (%.1f%%)", stats.CpuUsageMilli, stats.CpuLimit, stats.CpuUsagePercent)).SetTextColor(color)
		}
		return tview.NewTableCell(fmt.Sprintf("%.0fmCPU", stats.CpuUsageMilli))
	case 4:
		if stats.MemoryBytes <= 0 {
			return tview.NewTableCell("\u23F1").SetAlign(tview.AlignCenter)
		}

		//convet bytes to MiB
		memoryInMiB := stats.MemoryBytes / 1024 / 1024
		memoryLimitInMib := stats.MemoryLimitBytes / 1024 / 1024
		if stats.MemoryUsagePercent > 0 {
			color := t.getCellColor(stats.MemoryUsagePercent)
			return tview.NewTableCell(fmt.Sprintf("%.0f/%.0fMiB (%.1f%%)", memoryInMiB, memoryLimitInMib, stats.MemoryUsagePercent)).SetTextColor(color)
		}
		return tview.NewTableCell(fmt.Sprintf("%.0fMiB/-", memoryInMiB))
	default:
		return nil
	}
}

func (t *Table) getCellColor(utilization float64) tcell.Color {
	if utilization > 90 {
		return tcell.ColorRed
	}

	if utilization > 80 {
		return tcell.ColorYellow
	}

	return tcell.ColorWhite
}
