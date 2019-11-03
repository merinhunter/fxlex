package fxlex

import (
	"fmt"
	"io"
	"log"
	"testing"
)

var testFile = "test/lang.fx"

func TestNewLexer(t *testing.T) {
	if _, err := newLexer(testFile); err != nil {
		t.FailNow()
	}
}

func TestNewLexerError(t *testing.T) {
	if _, err := newLexer(""); err != nil {
	} else {
		t.FailNow()
	}
}

/*func TestLex(t *testing.T) {
	l, err := newLexer(testFile)
	if err != nil {
		t.FailNow()
	}

	for {
		if r, err := l.Lex(); err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		} else {
			fmt.Printf("%s:%d: %q\n", l.GetFilename(), l.GetLineNumber(), string(r))
		}
	}
}*/

func TestLex(t *testing.T) {
	l, err := newLexer(testFile)
	if err != nil {
		t.FailNow()
	}

	for t, err := l.Lex(); ; t, err = l.Lex() {
		if err != nil {
			if err == io.EOF {
				break
			} else {
				log.Fatal(err)
			}
		}

		fmt.Printf("%s:%d: (\"%v\";%v;%v)\n", l.GetFilename(), l.GetLineNumber(), t.lexeme, t.tokType, t.value)
	}
}
