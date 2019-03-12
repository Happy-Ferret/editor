package widget

import (
	"image"
	"image/color"

	"github.com/jmigpin/editor/util/drawutil"
	"github.com/jmigpin/editor/util/drawutil/drawer4"
	"github.com/jmigpin/editor/util/imageutil"
	"github.com/jmigpin/editor/util/iout/iorw"
)

type Text struct {
	ENode
	TextScroll

	Drawer   drawutil.Drawer
	OnSetStr func()

	scrollable struct{ x, y bool }
	ctx        ImageContext
	bg         color.Color

	brw iorw.ReadWriter // base rw
}

func NewText(ctx ImageContext) *Text {
	t := &Text{ctx: ctx}

	t.Drawer = drawer4.New()

	t.TextScroll.Text = t
	t.TextScroll.Drawer = t.Drawer

	t.brw = iorw.NewBytesReadWriter(nil)
	t.Drawer.SetReader(t.brw)

	return t
}

//----------

func (t *Text) Len() int {
	return t.brw.Len()
}

// Result might not be a copy, so changes to the slice might affect the text data.
func (t *Text) Bytes() ([]byte, error) {
	return t.brw.ReadNSliceAt(0, t.brw.Len())
}

func (t *Text) SetBytes(b []byte) error {
	if err := t.brw.Delete(0, t.brw.Len()); err != nil {
		return err
	}

	// run changes only once for delete+insert
	defer t.contentChanged()

	return t.brw.Insert(0, b)
}

//----------

func (t *Text) Str() string {
	p, err := t.Bytes()
	if err != nil {
		return ""
	}
	return string(p)
}

func (t *Text) SetStr(str string) error {
	return t.SetBytes([]byte(str))
}

//----------

func (t *Text) contentChanged() {
	t.Drawer.ContentChanged()

	// content changing can influence the layout in the case of dynamic sized textareas (needs layout). Also in the case of scrollareas that need to recalc scrollbars.
	t.MarkNeedsLayoutAndPaint()

	if t.OnSetStr != nil {
		t.OnSetStr()
	}
}

//----------

// implements Scrollable interface.
func (t *Text) SetScrollable(x, y bool) {
	t.scrollable.x = x
	t.scrollable.y = y
}

//----------

func (t *Text) RuneOffset() int {
	return t.Drawer.RuneOffset()
}

func (t *Text) SetRuneOffset(v int) {
	if t.scrollable.y && t.Drawer.RuneOffset() != v {
		t.Drawer.SetRuneOffset(v)
		t.MarkNeedsLayoutAndPaint()
	}
}

//----------

func (t *Text) IndexVisible(offset int) bool {
	return t.Drawer.RangeVisible(offset, 0)
}
func (t *Text) MakeIndexVisible(offset int) {
	t.MakeRangeVisible(offset, 0)
}
func (t *Text) MakeRangeVisible(offset, n int) {
	o := t.Drawer.RangeVisibleOffset(offset, n)
	t.SetRuneOffset(o)
}

//func (t *Text) MakeRangeVisibleCentered(offset, n int) {
//	o := t.Drawer.RangeVisibleOffsetCentered(offset, n)
//	t.SetRuneOffset(o)
//}

//----------

func (t *Text) GetPoint(i int) image.Point {
	return t.Drawer.PointOf(i)
}
func (t *Text) GetIndex(p image.Point) int {
	return t.Drawer.IndexOf(p)
}

//----------

func (t *Text) LineHeight() int {
	return t.Drawer.LineHeight()
}

//----------

func (t *Text) Measure(hint image.Point) image.Point {
	b := t.Bounds
	b.Max = b.Min.Add(hint)
	t.Drawer.SetBounds(b)
	m := t.Drawer.Measure()
	return imageutil.MinPoint(m, hint)
}

//----------

func (t *Text) Layout() {
	if t.Bounds != t.Drawer.Bounds() {
		t.Drawer.SetBounds(t.Bounds)
		t.MarkNeedsPaint()
	}
}

//----------

func (t *Text) PaintBase() {
	imageutil.FillRectangle(t.ctx.Image(), &t.Bounds, t.bg)
}
func (t *Text) Paint() {
	t.Drawer.Draw(t.ctx.Image())
}

//----------

func (t *Text) OnThemeChange() {
	fg := t.TreeThemePaletteColor("text_fg")
	t.Drawer.SetFg(fg)

	t.bg = t.TreeThemePaletteColor("text_bg")

	f := t.TreeThemeFont().Face(nil)
	if f != t.Drawer.Face() {
		t.Drawer.SetFace(f)
		t.MarkNeedsLayoutAndPaint()
	} else {
		t.MarkNeedsPaint()
	}
}
