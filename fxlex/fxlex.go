package fxlex

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"unicode"
)

type Token struct {
	lexeme  string
	tokType TokType
	value   int64
}

type TokType int

const (
	RuneEOF         = 0
	TokKey  TokType = iota
	TokId
	TokFunc

	// Basic Types
	TokInt
	TokIntLit
	TokBool
	TokBoolLit
	TokCoord

	// Punctuation tokens
	TokLPar     TokType = '('
	TokRPar     TokType = ')'
	TokLCurl    TokType = '{'
	TokRCurl    TokType = '}'
	TokLSquare  TokType = '['
	TokRSquare  TokType = ']'
	TokComma    TokType = ','
	TokDot      TokType = '.'
	Semicolon   TokType = ';'
	Assignation TokType = '='
	Declaration TokType = iota

	// Int operators
	TokPlus   TokType = '+'
	TokMinus  TokType = '-'
	TokTimes  TokType = '*'
	TokDivide TokType = '/'
	TokRem    TokType = '%'
	TokGT     TokType = '>'
	TokLT     TokType = '<'
	TokPow    TokType = iota
	TokGTE
	TokLTE

	// Bool operators
	TokOr  TokType = '|'
	TokAnd TokType = '&'
	TokNeg TokType = '!'
	TokXor TokType = '^'
)

type RuneScanner interface {
	ReadRune() (r rune, size int, err error)
	UnreadRune() error
}

type Lexer struct {
	file     string
	fd       *os.File
	line     int
	rs       RuneScanner
	lastrune rune

	accepted []rune
	tokSaved *Token
}

func newLexer(file string) (l *Lexer, err error) {
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	l = &Lexer{file: file, fd: fd, line: 1, rs: bufio.NewReader(fd)}
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
	if l.lastrune == RuneEOF {
		return
	}

	err := l.rs.UnreadRune()
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
			} else {
				return t, errors.New("bad declaration token")
			}

		// Punctuation tokens
		case '(', ')', '{', '}', '[', ']', ',', '.', ';', '=':
			t.tokType = TokType(r)
			t.lexeme = l.accept()
			return t, nil

		// Int operators
		case '+', '-', '*', '/', '%', '>', '<':
			t.tokType = TokType(r)

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
			t.tokType = TokType(r)
			t.lexeme = l.accept()
			return t, nil

		// EOF
		case RuneEOF:
			err := l.fd.Close()
			if err != nil {
				panic(err)
			}

			l.accept()
			return t, io.EOF
		}

		switch {
		case unicode.IsDigit(r):
			l.unget()
			t, err = l.lexNum()
			return t, err
		case unicode.IsLetter(r):
			l.unget()
			t, err = l.lexId()

			isKeyword := func(lexeme string) bool {
				keywords := map[string]bool{
					"type":   true,
					"record": true,
					"circle": true,
					"rect":   true,
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

				if t.lexeme == "int" {
					t.tokType = TokInt
				}

				if t.lexeme == "bool" {
					t.tokType = TokBool
				}

				if t.lexeme == "Coord" {
					t.tokType = TokCoord
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

	return t, err
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

func (l *Lexer) lexId() (t Token, err error) {
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

	t.tokType = TokId
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
