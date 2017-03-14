package tautil

func DuplicateLines(ta Texta) {
	a, b, hasNewline := linesStringIndexes(ta)
	t := ta.Str()[a:b]
	if !hasNewline {
		ta.EditInsert(b, "\n")
		b++
	}
	ta.EditInsert(b, t)
	ta.EditDone()
	ta.SetSelectionOn(true)
	ta.SetSelectionIndex(b)

	// cursor index without the newline
	c := len(t)
	if hasNewline {
		_, u, ok := PreviousRuneIndex(t, c)
		if ok {
			c = u
		}
	}
	ta.SetCursorIndex(b + c)
}