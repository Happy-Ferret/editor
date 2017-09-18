package toolbardata

import (
	"strconv"
	"strings"
	"unicode"
)

type ToolbarData struct {
	Str   string
	Parts []*Part
}

func NewToolbarData(str string) *ToolbarData {
	td := &ToolbarData{Str: str}
	td.Parts = parseParts(str)
	return td
}

func (td *ToolbarData) GetPartAtIndex(i int) (*Part, bool) {
	for _, p := range td.Parts {
		if i >= p.S && i < p.E {
			return p, true
		}
	}
	// return last part for index at eos
	if i == len(td.Str) && len(td.Parts) > 0 {
		return td.Parts[len(td.Parts)-1], true
	}
	return nil, false
}

func (td *ToolbarData) ReplacePart(i int, str string) string {

	p := td.Parts[i]
	return td.Str[:p.S] + str + td.Str[p.E:]
}

func (td *ToolbarData) part0Arg0Token() (*Token, bool) {
	if len(td.Parts) == 0 {
		return nil, false
	}
	if len(td.Parts[0].Args) == 0 {
		return nil, false
	}
	return td.Parts[0].Args[0], true
}

func (td *ToolbarData) DecodePart0Arg0() string {
	tok, ok := td.part0Arg0Token()
	if !ok {
		return ""
	}
	return RemoveHomeVars(tok.Str)
}

func (td *ToolbarData) StrWithPart0Arg0Encoded() string {
	tok, ok := td.part0Arg0Token()
	if !ok {
		return td.Str
	}
	s2 := RemoveHomeVars(tok.Str)
	s3 := InsertHomeVars(s2)
	return td.Str[:tok.S] + s3 + td.Str[tok.E:]
}

func (td *ToolbarData) StrWithPart0Arg0Decoded() string {
	tok, ok := td.part0Arg0Token()
	if !ok {
		return td.Str
	}
	s2 := RemoveHomeVars(tok.Str)
	return td.Str[:tok.S] + s2 + td.Str[tok.E:]
}

func parseParts(str string) []*Part {
	var parts []*Part
	toks := parseTokens(str, 0, len(str), '|')
	for _, t := range toks {
		ctoks := parseTokens(str, t.S, t.E, ' ')
		ctoks = filterEmptyTokens(ctoks)
		p := &Part{Token: *t, Args: ctoks}
		parts = append(parts, p)
	}
	return parts
}
func parseTokens(str string, a, b int, sep rune) []*Token {
	lastQuote := rune(0)
	escape := false
	split := func(ru rune) bool {
		switch {
		case ru == '\\':
			escape = true
			return false
		case escape:
			escape = false
			return false
		case ru == lastQuote:
			lastQuote = 0
			return false
		case lastQuote != 0: // inside a quote
			return false
		case unicode.In(ru, unicode.Quotation_Mark):
			lastQuote = ru
			return false
		default:
			return ru == sep
		}
	}
	return fieldsFunc(str, a, b, split)
}
func fieldsFunc(str string, a, b int, split func(rune) bool) []*Token {
	var u []*Token
	s := a
	for i, ru := range str[a:b] {
		if split(ru) {
			t := NewToken(str, s, a+i)
			s = a + i + len(string(ru)) // not including separator in tok
			u = append(u, t)
		}
	}
	if s < b {
		t := NewToken(str, s, b)
		u = append(u, t)
	}
	return u
}
func filterEmptyTokens(toks []*Token) []*Token {
	var u []*Token
	for _, t := range toks {
		if !t.isEmpty() {
			u = append(u, t)
		}
	}
	return u
}

type Part struct {
	Token
	Args []*Token
}

type Token struct {
	Str  string // token string
	S, E int    // start/end str indexes of the root string
}

func NewToken(str string, s, e int) *Token {
	tok := &Token{Str: str[s:e], S: s, E: e}

	// unquote str if possible
	str2, err := strconv.Unquote(tok.Str)
	if err == nil {
		v, _, _, err2 := strconv.UnquoteChar(tok.Str, 0)
		if err2 != nil {
			panic(err2)
		}

		l := len(string(v))
		tok.S += l
		tok.E -= l
		tok.Str = str2
	}

	return tok
}

func (tok *Token) isEmpty() bool {
	return strings.TrimSpace(tok.Str) == ""
}
