package fxlex

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

const (
	RuneEOF = unicode.MaxRune + 1 + iota
	TokKey
	TokID
	TokFunc

	// Basic Types
	TokIntLit
	TokBoolLit

	// Punctuation tokens
	Declaration

	// Int operators
	TokPow
	TokGTE
	TokLTE
)

const (
	// Punctuation tokens
	TokLPar     = '('
	TokRPar     = ')'
	TokLCurl    = '{'
	TokRCurl    = '}'
	TokLSquare  = '['
	TokRSquare  = ']'
	TokComma    = ','
	TokDot      = '.'
	Semicolon   = ';'
	Assignation = '='

	// Int operators
	TokPlus   = '+'
	TokMinus  = '-'
	TokTimes  = '*'
	TokDivide = '/'
	TokRem    = '%'
	TokGT     = '>'
	TokLT     = '<'

	// Bool operators
	TokOr  = '|'
	TokAnd = '&'
	TokNeg = '!'
	TokXor = '^'

	TokEOF = RuneEOF
)

type tokType int

type RuneScanner interface {
	ReadRune() (r rune, size int, err error)
	UnreadRune() error
}

type Token struct {
	lexeme  string
	tokType int
	value   int64
}

type Lexer struct {
	file     string
	line     int
	rs       RuneScanner
	lastrune rune

	accepted []rune
	tokSaved *Token
}

func NewLexer(rs RuneScanner, file string) (l *Lexer, err error) {
	l = &Lexer{line: 1}
	l.file = file
	l.rs = rs

	return l, nil
}

func (l *Lexer) get() (r rune) {
	var err error
	r, _, err = l.rs.ReadRune()

	if err == nil {
		l.lastrune = r
		if r == '\n' {
			l.line++
		}
	}

	if err == io.EOF {
		l.lastrune = RuneEOF
		return RuneEOF
	}

	// Panic for unexpected errors
	if err != nil {
		panic(err)
	}

	l.accepted = append(l.accepted, r)
	return r
}

func (l *Lexer) unget() {
	var err error
	if l.lastrune == RuneEOF {
		return
	}

	err = l.rs.UnreadRune()
	if err == nil && l.lastrune == '\n' {
		l.line--
	}

	l.lastrune = unicode.ReplacementChar

	if len(l.accepted) != 0 {
		l.accepted = l.accepted[0 : len(l.accepted)-1]
	}

	if err != nil {
		panic(err)
	}
}

func (l *Lexer) accept() (tok string) {
	tok = string(l.accepted)

	if tok == "" && l.lastrune != RuneEOF {
		panic(errors.New("empty token"))
	}

	l.accepted = nil
	return tok
}

func (l *Lexer) skipComment() {
	isOver := func(ar rune) bool {
		return ar == '\n' || ar == RuneEOF
	}

	for r := l.get(); !isOver(r); r = l.get() {
	}

	l.accept()
}

func (l *Lexer) Lex() (t Token, err error) {
	if l.tokSaved != nil {
		t = *l.tokSaved
		l.tokSaved = nil
		return t, nil
	}

	for r := l.get(); ; r = l.get() {
		if unicode.IsSpace(r) {
			l.accept()
			continue
		}

		switch r {
		// Declaration
		case ':':
			if l.get() == '=' {
				t.tokType = Declaration
				t.lexeme = l.accept()
				return t, nil
			}
			return t, errors.New("bad declaration token")

		// Punctuation tokens
		case '(', ')', '{', '}', '[', ']', ',', '.', ';', '=':
			t.tokType = int(r)
			t.lexeme = l.accept()
			return t, nil

		// Int operators
		case '+', '-', '*', '/', '%', '>', '<':
			t.tokType = int(r)

			if l.get() == '/' && r == '/' {
				l.skipComment()
				continue
			} else {
				l.unget()
			}

			if l.get() == '*' && r == '*' {
				t.tokType = TokPow
			} else {
				l.unget()
			}

			if l.get() == '=' && r == '>' {
				t.tokType = TokGTE
			} else {
				l.unget()
			}

			if l.get() == '=' && r == '<' {
				t.tokType = TokLTE
			} else {
				l.unget()
			}

			t.lexeme = l.accept()
			return t, nil

		// Bool operators
		case '|', '&', '!', '^':
			t.tokType = int(r)
			t.lexeme = l.accept()
			return t, nil

		// EOF
		case RuneEOF:
			t.tokType = TokEOF
			l.accept()
			return t, nil
		}

		switch {
		case unicode.IsDigit(r):
			l.unget()
			t, err = l.lexNum()
			return t, err

		case unicode.IsLetter(r):
			l.unget()
			t, err = l.lexID()

			isKeyword := func(lexeme string) bool {
				keywords := map[string]bool{
					"type":   true,
					"record": true,
					"iter":   true,
					"if":     true,
					"else":   true,
				}

				return keywords[lexeme]
			}

			if err == nil {
				if isKeyword(t.lexeme) {
					t.tokType = TokKey
				}

				if t.lexeme == "func" {
					t.tokType = TokFunc
				}

				if t.lexeme == "True" {
					t.tokType = TokBoolLit
					t.value = 1
				}

				if t.lexeme == "False" {
					t.tokType = TokBoolLit
					t.value = 0
				}
			}
			return t, err

		default:
			errs := fmt.Sprintf("bad rune %c: %x", r, r)
			return t, errors.New(errs)
		}
	}
}

func (l *Lexer) lexNum() (t Token, err error) {
	hex := false
	base := 0
	bitSize := 64

	r := l.get()
	if l.get() == 'x' && r == '0' {
		hex = true
	} else {
		l.unget()
	}

	isHex := func(ar rune) bool {
		return unicode.IsDigit(ar) || (ar >= 'a' && ar <= 'f') || (ar >= 'A' && ar <= 'F')
	}

	if hex {
		for r = l.get(); isHex(r); r = l.get() {
		}
	} else {
		for r = l.get(); unicode.IsDigit(r); r = l.get() {
		}
	}
	l.unget()

	t.lexeme = l.accept()
	t.tokType = TokIntLit
	t.value, err = strconv.ParseInt(t.lexeme, base, bitSize)

	if err != nil {
		return t, errors.New("bad int [" + t.lexeme + "]")
	}

	return t, nil
}

func (l *Lexer) lexID() (t Token, err error) {
	r := l.get()

	if !unicode.IsLetter(r) {
		return t, errors.New("bad ID")
	}

	isAlpha := func(ar rune) bool {
		return unicode.IsDigit(ar) || unicode.IsLetter(ar) || r == '_'
	}

	for r = l.get(); isAlpha(r); r = l.get() {
	}
	l.unget()

	t.tokType = TokID
	t.lexeme = l.accept()
	return t, nil
}

func (l *Lexer) Peek() (t Token, err error) {
	t, err = l.Lex()
	if err == nil {
		l.tokSaved = &t
	}
	return t, nil
}

func (l *Lexer) GetFilename() string {
	return l.file
}

func (l *Lexer) GetLineNumber() int {
	return l.line
}

func (t Token) String() string {
	return fmt.Sprintf("{\"%s\",%s,%d}", t.lexeme, tokType(t.tokType), t.value)
}

func (t Token) GetLexeme() string {
	return t.lexeme
}

func (t Token) GetTokType() int {
	return t.tokType
}

func (t Token) GetValue() int64 {
	return t.value
}
