package ui

import (
	"image"
	"image/color"

	"github.com/jmigpin/editor/uiutil/widget"
)

type Columns struct {
	widget.EndPercentLayout
	Layout *Layout
}

func NewColumns(layout *Layout) *Columns {
	cols := &Columns{Layout: layout}

	cols.NewColumn() // start with 1 column

	return cols
}
func (cols *Columns) Paint() {
	if len(cols.Childs()) == 0 {
		b := cols.Bounds()
		cols.Layout.UI.FillRectangle(&b, color.White)
	}
}
func (cols *Columns) LastColumnOrNew() *Column {
	col, ok := cols.LastChildColumn()
	if !ok {
		col = cols.NewColumn()
	}
	return col
}
func (cols *Columns) NewColumn() *Column {
	col := NewColumn(cols)
	cols.insertColumnBefore(col, nil)
	return col
}
func (cols *Columns) insertColumnBefore(col, next *Column) {
	if next == nil {

		// TODO: need to return false

		//// don't insert if it will be too small
		//lc := cols.LastChild()
		//if lc != nil && lc.Prev() != nil {
		//	start := cols.ChildEndPercent(lc.Prev())
		//	end := cols.ChildEndPercent(lc)
		//	x := int((end - start) * float64(cols.Bounds().Dx()))
		//	if x < 40 {
		//		return
		//	}
		//}

		widget.PushBack(cols, col)
	} else {
		panic("TODO")
		widget.InsertBefore(cols, col, next)
	}

	cols.fixFirstColSeparator()
	cols.CalcChildsBounds()
	cols.MarkNeedsPaint()
}

// TODO: override node.Remove?

func (cols *Columns) removeColumn(col *Column) {
	cols.Remove(col)
	cols.fixFirstColSeparator()
	cols.CalcChildsBounds()
	cols.MarkNeedsPaint()
}

func (cols *Columns) fixFirstColSeparator() {
	for i, c := range cols.Columns() {
		c.HideSeparator(i == 0)
	}
}

func (cols *Columns) CloseColumnEnsureOne(col *Column) {
	col.Close()
	// ensure one column
	if len(cols.Childs()) == 0 {
		_ = cols.NewColumn()
	}
}

// Used by restore session.
func (cols *Columns) CloseAllAndOpenN(n int) {
	// close all columns
	for len(cols.Childs()) > 0 {
		u, _ := cols.FirstChildColumn()
		u.Close()
	}
	// ensure one column
	if n <= 1 {
		n = 1
	}
	// n new columns
	for i := 0; i < n; i++ {
		_ = cols.NewColumn()
	}
}

func (cols *Columns) PointColumn(p *image.Point) (*Column, bool) {
	for _, c := range cols.Columns() {
		if p.In(c.Bounds()) {
			return c, true
		}
	}
	return nil, false
}

func (cols *Columns) FirstChildColumn() (*Column, bool) {
	u := cols.FirstChild()
	if u == nil {
		return nil, false
	}
	return u.(*Column), true
}
func (cols *Columns) LastChildColumn() (*Column, bool) {
	u := cols.LastChild()
	if u == nil {
		return nil, false
	}
	return u.(*Column), true
}
func (cols *Columns) Columns() []*Column {
	childs := cols.Childs()
	u := make([]*Column, 0, len(childs))
	for _, h := range childs {
		u = append(u, h.(*Column))
	}
	return u
}
