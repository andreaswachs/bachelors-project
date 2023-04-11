package pp

import (
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
)

type RowData interface {
	any
}

type PrettyPrinter interface {
	AddHeader(header ...string)
	AddRow(row ...RowData)
	Print()
}

type prettyPrinter struct {
	header []string
	rows   [][]RowData

	shouldPrint bool

	table table.Writer
}

func NewPrettyPrinter() PrettyPrinter {
	return &prettyPrinter{
		header: make([]string, 0),
		rows:   make([][]RowData, 0),
		table:  table.NewWriter(),
	}
}

func (pp *prettyPrinter) AddHeader(header ...string) {
	pp.header = header
	pp.shouldPrint = true
}

func (pp *prettyPrinter) AddRow(data ...RowData) {
	pp.rows = append(pp.rows, data)
}

func (pp *prettyPrinter) Print() {
	if !pp.shouldPrint {
		return
	}

	pp.table.SetOutputMirror(os.Stdout)
	pp.table.SetStyle(table.StyleLight)

	headerRow := make([]interface{}, 0)
	for _, data := range pp.header {
		headerRow = append(headerRow, data)
	}
	pp.table.AppendHeader(headerRow)

	for _, row := range pp.rows {
		newRow := make([]interface{}, 0)
		for _, data := range row {
			newRow = append(newRow, data)
		}

		pp.table.AppendRow(newRow)
	}

	pp.table.Render()
}
