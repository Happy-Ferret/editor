package ui

import (
	"image"
	"image/color"

	// only for cursordef

	"github.com/jmigpin/editor/uiutil/widget"
	"github.com/jmigpin/editor/xgbutil/evreg"
	"github.com/jmigpin/editor/xgbutil/xinput"
)

// Used in row and column to move and close.
type Square struct {
	widget.EmbedNode
	ui      *UI
	EvReg   *evreg.Register
	evUnreg evreg.Unregister

	pressPointPad image.Point
	buttonPressed bool
	values        [7]bool // bg and mini-squares

	Width int

	ColumnStyle bool
}

func NewSquare(ui *UI) *Square {
	sq := &Square{ui: ui}
	sq.Width = SquareWidth

	//width := SquareWidth
	//sq.C.Style.MainSize = &width

	sq.EvReg = evreg.NewRegister()
	r1 := sq.ui.EvReg.Add(xinput.ButtonPressEventId,
		&evreg.Callback{sq.onButtonPress})
	r2 := sq.ui.EvReg.Add(xinput.ButtonReleaseEventId,
		&evreg.Callback{sq.onButtonRelease})
	r3 := sq.ui.EvReg.Add(xinput.MotionNotifyEventId,
		&evreg.Callback{sq.onMotionNotify})
	sq.evUnreg.Add(r1, r2, r3)

	//sq.ui.CursorMan.SetBoundsCursor(&sq.Bounds(), xcursor.Icon)

	return sq
}
func (sq *Square) Close() {
	sq.evUnreg.UnregisterAll()
	//sq.ui.CursorMan.RemoveBoundsCursor(&sq.Bounds())
}

func (sq *Square) Measure(hint image.Point) image.Point {
	return image.Point{sq.Width, sq.Width}
}
func (sq *Square) CalcChildsBounds() {
}
func (sq *Square) Paint() {
	var c color.Color = SquareColor
	if sq.values[SquareEdited] {
		c = SquareEditedColor
	}
	if sq.values[SquareNotExist] {
		c = SquareNotExistColor
	}
	if sq.values[SquareExecuting] {
		c = SquareExecutingColor
	}
	bounds := sq.Bounds()
	sq.ui.FillRectangle(&bounds, c)

	if sq.values[SquareDuplicate] {
		c2 := SquareEditedColor
		sq.ui.BorderRectangle(&bounds, c2, 2)
	}

	// separator
	r3 := bounds
	if ScrollbarLeft {
		r3.Min.X = r3.Max.X - 1
	} else {
		r3.Max.X = r3.Min.X + 1
	}
	sq.ui.FillRectangle(&r3, RowInnerSeparatorColor)

	if sq.ColumnStyle {
		r4 := bounds
		r4.Min.Y = r4.Max.Y - 1
		sq.ui.FillRectangle(&r4, RowInnerSeparatorColor)
	}

	// mini-squares

	miniSq := func(i int) *image.Rectangle {
		r := bounds
		w := (r.Max.X - r.Min.X) / 2
		r.Max.X = r.Min.X + w
		r.Max.Y = r.Min.Y + w
		switch i {
		case 0:
		case 1:
			r.Min.X += w
			r.Max.X += w
		case 2:
			r.Min.Y += w
			r.Max.Y += w
		case 3:
			r.Min.X += w
			r.Max.X += w
			r.Min.Y += w
			r.Max.Y += w
		}
		return &r
	}

	if sq.values[SquareDiskChanges] {
		u := 0
		if ScrollbarLeft {
			u = 1
		}
		r := miniSq(u)
		sq.ui.FillRectangle(r, SquareDiskChangesColor)
	}
	if sq.values[SquareActive] {
		u := 1
		if ScrollbarLeft {
			u = 0
		}
		r := miniSq(u)
		sq.ui.FillRectangle(r, SquareActiveColor)
	}
}
func (sq *Square) onButtonPress(ev0 interface{}) {
	if sq.Hidden() {
		return
	}

	ev := ev0.(*xinput.ButtonPressEvent)
	if !ev.Point.In(sq.Bounds()) {
		return
	}
	sq.buttonPressed = true

	var u image.Point
	if ScrollbarLeft {
		u = image.Point{sq.Bounds().Min.X, sq.Bounds().Min.Y}
	} else {
		u = image.Point{sq.Bounds().Max.X, sq.Bounds().Min.Y}
	}
	sq.pressPointPad = u.Sub(*ev.Point)

	ev2 := &SquareButtonPressEvent{sq, ev.Button, ev.Point}
	sq.EvReg.RunCallbacks(SquareButtonPressEventId, ev2)
}
func (sq *Square) onButtonRelease(ev0 interface{}) {
	if !sq.buttonPressed {
		return
	}
	sq.buttonPressed = false
	ev := ev0.(*xinput.ButtonReleaseEvent)
	ev2 := &SquareButtonReleaseEvent{sq, ev.Button, ev.Point}
	sq.EvReg.RunCallbacks(SquareButtonReleaseEventId, ev2)
}
func (sq *Square) onMotionNotify(ev0 interface{}) {
	if !sq.buttonPressed {
		return
	}
	ev := ev0.(*xinput.MotionNotifyEvent)
	ev2 := &SquareMotionNotifyEvent{sq, ev.Mods, ev.Point, &sq.pressPointPad}
	sq.EvReg.RunCallbacks(SquareMotionNotifyEventId, ev2)
}

func (sq *Square) WarpPointer() {
	sa := sq.Bounds()
	p := sa.Min.Add(image.Pt(sa.Dx()/2, sa.Dy()/2))
	sq.ui.WarpPointer(&p)
}

func (sq *Square) Value(t SquareType) bool {
	return sq.values[t]
}
func (sq *Square) SetValue(t SquareType, v bool) {
	if sq.values[t] != v {
		sq.values[t] = v
		sq.MarkNeedsPaint()
	}
}

type SquareType int

const (
	SquareNone SquareType = iota
	SquareActive
	SquareExecuting
	SquareEdited
	SquareDiskChanges
	SquareNotExist
	SquareDuplicate
)

const (
	SquareButtonPressEventId = iota
	SquareButtonReleaseEventId
	SquareMotionNotifyEventId
)

type SquareButtonPressEvent struct {
	Square *Square
	Button *xinput.Button
	Point  *image.Point
}
type SquareButtonReleaseEvent struct {
	Square *Square
	Button *xinput.Button
	Point  *image.Point
}
type SquareMotionNotifyEvent struct {
	Square *Square
	Mods   xinput.Modifiers
	Point  *image.Point

	PressPointPad *image.Point
}
