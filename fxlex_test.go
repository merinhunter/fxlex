package fxlex_test

import (
	"bufio"
	"fmt"
	. "fxlex"
	"log"
	"strings"
	"testing"
)

type testToken struct {
	lexeme  string
	tokType int
	value   int64
}

func (t testToken) String() string {
	return fmt.Sprintf("{\"%s\",toktype(%d),%d}", t.lexeme, t.tokType, t.value)
}

var examplesList = []string{
	`type record vector(int x, int y, int z)`,
	`//comment`,
	`p.x = v.x*i;`,
	`iter (i := 0, v.z, 2){	//declares the variable only in the loop`,
	`circle(pp, 0x3, 0x1100001f);`,
	`rect(pp, α, 0xff);`,
	`p.y = v.y**i;`,
	`bool check = False;`,
}

var tokensList = [][]testToken{
	{
		{"type", TokKey, 0},
		{"record", TokKey, 0},
		{"vector", TokID, 0},
		{"(", TokLPar, 0},
		{"int", TokID, 0},
		{"x", TokID, 0},
		{",", TokComma, 0},
		{"int", TokID, 0},
		{"y", TokID, 0},
		{",", TokComma, 0},
		{"int", TokID, 0},
		{"z", TokID, 0},
		{")", TokRPar, 0},
	}, {}, {
		{"p", TokID, 0},
		{".", TokDot, 0},
		{"x", TokID, 0},
		{"=", Assignation, 0},
		{"v", TokID, 0},
		{".", TokDot, 0},
		{"x", TokID, 0},
		{"*", TokTimes, 0},
		{"i", TokID, 0},
		{";", Semicolon, 0},
	}, {
		{"iter", TokKey, 0},
		{"(", TokLPar, 0},
		{"i", TokID, 0},
		{":=", Declaration, 0},
		{"0", TokIntLit, 0},
		{",", TokComma, 0},
		{"v", TokID, 0},
		{".", TokDot, 0},
		{"z", TokID, 0},
		{",", TokComma, 0},
		{"2", TokIntLit, 2},
		{")", TokRPar, 0},
		{"{", TokLCurl, 0},
	}, {
		{"circle", TokID, 0},
		{"(", TokLPar, 0},
		{"pp", TokID, 0},
		{",", TokComma, 0},
		{"0x3", TokIntLit, 3},
		{",", TokComma, 0},
		{"0x1100001f", TokIntLit, 285212703},
		{")", TokRPar, 0},
		{";", Semicolon, 0},
	}, {
		{"rect", TokID, 0},
		{"(", TokLPar, 0},
		{"pp", TokID, 0},
		{",", TokComma, 0},
		{"α", TokID, 0},
		{",", TokComma, 0},
		{"0xff", TokIntLit, 255},
		{")", TokRPar, 0},
		{";", Semicolon, 0},
	}, {
		{"p", TokID, 0},
		{".", TokDot, 0},
		{"y", TokID, 0},
		{"=", Assignation, 0},
		{"v", TokID, 0},
		{".", TokDot, 0},
		{"y", TokID, 0},
		{"**", TokPow, 0},
		{"i", TokID, 0},
		{";", Semicolon, 0},
	}, {
		{"bool", TokID, 0},
		{"check", TokID, 0},
		{"=", Assignation, 0},
		{"False", TokBoolLit, 0},
		{";", Semicolon, 0},
	},
}

var exampleFile = `//basic types bool, int (64 bits), Coord(int x, int y)
//literals of type int are 2, 3, or 0x2dfadfd
//literals of Coord are [3,4] [0x46,4]
//literals of bool are True, False
//operators of int are + - * / ** > >= < <=
//operators of int are %
//operators of bool are | & ! ^
//precedence is like in C, with ** having the
//same precedence as sizeof (not present in fx)

type record vector(int x, int y, int z)
type record difficult (vector v, Coord r)

//builtins
//circle(p, 2, 0x1100001f);
//	at point p, int radius r, color: transparency and rgb
//rect(p, α, col);
//	at point p, int angle (degrees),
//	color: transparency (0-100) and rgb

//macro definition
func line(vector v){
	Coord p;		//only local variables, no globals

				//last number in the loop is the step
	iter (i := 0, v.z, 2){	//declares de variable only in the loop
		p.x = v.x*i;
		p.y = v.y*i;
		circle(p, 2, 1);
	}
}

//macro entry
func main(){
	vector v;
	Coord pp;

	v.x = 3;
	v.y = 8;
	v.z = 2;
	pp = [4,45];
	if(v.x > 3 | True) {		// (v.x>3)|True
		circle(pp, 0x3, 0x1100001f);
	} else {
		line(v);
		line(v);
	}
	line(v);
	line(v);
	iter (i := 0, 3, 1){		//loops 0 1 2 3
		rect(pp, α, 0xff);
	}
}
`

var tokensFile = []testToken{
	{"type", TokKey, 0},
	{"record", TokKey, 0},
	{"vector", TokID, 0},
	{"(", TokLPar, 0},
	{"int", TokID, 0},
	{"x", TokID, 0},
	{",", TokComma, 0},
	{"int", TokID, 0},
	{"y", TokID, 0},
	{",", TokComma, 0},
	{"int", TokID, 0},
	{"z", TokID, 0},
	{")", TokRPar, 0},
	{"type", TokKey, 0},
	{"record", TokKey, 0},
	{"difficult", TokID, 0},
	{"(", TokLPar, 0},
	{"vector", TokID, 0},
	{"v", TokID, 0},
	{",", TokComma, 0},
	{"Coord", TokID, 0},
	{"r", TokID, 0},
	{")", TokRPar, 0},
	{"func", TokFunc, 0},
	{"line", TokID, 0},
	{"(", TokLPar, 0},
	{"vector", TokID, 0},
	{"v", TokID, 0},
	{")", TokRPar, 0},
	{"{", TokLCurl, 0},
	{"Coord", TokID, 0},
	{"p", TokID, 0},
	{";", Semicolon, 0},
	{"iter", TokKey, 0},
	{"(", TokLPar, 0},
	{"i", TokID, 0},
	{":=", Declaration, 0},
	{"0", TokIntLit, 0},
	{",", TokComma, 0},
	{"v", TokID, 0},
	{".", TokDot, 0},
	{"z", TokID, 0},
	{",", TokComma, 0},
	{"2", TokIntLit, 2},
	{")", TokRPar, 0},
	{"{", TokLCurl, 0},
	{"p", TokID, 0},
	{".", TokDot, 0},
	{"x", TokID, 0},
	{"=", Assignation, 0},
	{"v", TokID, 0},
	{".", TokDot, 0},
	{"x", TokID, 0},
	{"*", TokTimes, 0},
	{"i", TokID, 0},
	{";", Semicolon, 0},
	{"p", TokID, 0},
	{".", TokDot, 0},
	{"y", TokID, 0},
	{"=", Assignation, 0},
	{"v", TokID, 0},
	{".", TokDot, 0},
	{"y", TokID, 0},
	{"*", TokTimes, 0},
	{"i", TokID, 0},
	{";", Semicolon, 0},
	{"circle", TokID, 0},
	{"(", TokLPar, 0},
	{"p", TokID, 0},
	{",", TokComma, 0},
	{"2", TokIntLit, 2},
	{",", TokComma, 0},
	{"1", TokIntLit, 1},
	{")", TokRPar, 0},
	{";", Semicolon, 0},
	{"}", TokRCurl, 0},
	{"}", TokRCurl, 0},
	{"func", TokFunc, 0},
	{"main", TokID, 0},
	{"(", TokLPar, 0},
	{")", TokRPar, 0},
	{"{", TokLCurl, 0},
	{"vector", TokID, 0},
	{"v", TokID, 0},
	{";", Semicolon, 0},
	{"Coord", TokID, 0},
	{"pp", TokID, 0},
	{";", Semicolon, 0},
	{"v", TokID, 0},
	{".", TokDot, 0},
	{"x", TokID, 0},
	{"=", Assignation, 0},
	{"3", TokIntLit, 3},
	{";", Semicolon, 0},
	{"v", TokID, 0},
	{".", TokDot, 0},
	{"y", TokID, 0},
	{"=", Assignation, 0},
	{"8", TokIntLit, 8},
	{";", Semicolon, 0},
	{"v", TokID, 0},
	{".", TokDot, 0},
	{"z", TokID, 0},
	{"=", Assignation, 0},
	{"2", TokIntLit, 2},
	{";", Semicolon, 0},
	{"pp", TokID, 0},
	{"=", Assignation, 0},
	{"[", TokLSquare, 0},
	{"4", TokIntLit, 4},
	{",", TokComma, 0},
	{"45", TokIntLit, 45},
	{"]", TokRSquare, 0},
	{";", Semicolon, 0},
	{"if", TokKey, 0},
	{"(", TokLPar, 0},
	{"v", TokID, 0},
	{".", TokDot, 0},
	{"x", TokID, 0},
	{">", TokGT, 0},
	{"3", TokIntLit, 3},
	{"|", TokOr, 0},
	{"True", TokBoolLit, 1},
	{")", TokRPar, 0},
	{"{", TokLCurl, 0},
	{"circle", TokID, 0},
	{"(", TokLPar, 0},
	{"pp", TokID, 0},
	{",", TokComma, 0},
	{"0x3", TokIntLit, 3},
	{",", TokComma, 0},
	{"0x1100001f", TokIntLit, 285212703},
	{")", TokRPar, 0},
	{";", Semicolon, 0},
	{"}", TokRCurl, 0},
	{"else", TokKey, 0},
	{"{", TokLCurl, 0},
	{"line", TokID, 0},
	{"(", TokLPar, 0},
	{"v", TokID, 0},
	{")", TokRPar, 0},
	{";", Semicolon, 0},
	{"line", TokID, 0},
	{"(", TokLPar, 0},
	{"v", TokID, 0},
	{")", TokRPar, 0},
	{";", Semicolon, 0},
	{"}", TokRCurl, 0},
	{"line", TokID, 0},
	{"(", TokLPar, 0},
	{"v", TokID, 0},
	{")", TokRPar, 0},
	{";", Semicolon, 0},
	{"line", TokID, 0},
	{"(", TokLPar, 0},
	{"v", TokID, 0},
	{")", TokRPar, 0},
	{";", Semicolon, 0},
	{"iter", TokKey, 0},
	{"(", TokLPar, 0},
	{"i", TokID, 0},
	{":=", Declaration, 0},
	{"0", TokIntLit, 0},
	{",", TokComma, 0},
	{"3", TokIntLit, 3},
	{",", TokComma, 0},
	{"1", TokIntLit, 1},
	{")", TokRPar, 0},
	{"{", TokLCurl, 0},
	{"rect", TokID, 0},
	{"(", TokLPar, 0},
	{"pp", TokID, 0},
	{",", TokComma, 0},
	{"α", TokID, 0},
	{",", TokComma, 0},
	{"0xff", TokIntLit, 255},
	{")", TokRPar, 0},
	{";", Semicolon, 0},
	{"}", TokRCurl, 0},
	{"}", TokRCurl, 0},
}

var example2 = `//macro entry
func main(){
	vector v;
	Coord pp;

	v.x = 3;
	v.y = 8;
	v.z := 2;
	pp = [4,45];
	if(v.x > 3 | True) {		// (v.x>3)|True
		circle(pp, 0x3, 0x1100001f);
	} else {
		line(v);
		line(v);
	}
	line(v);
	line(v);
	iter (i := 0, 3, 1){		//loops 0 1 2 3
		rect(pp, α, 0xff);
	}
}`

func newTestLexer(t *testing.T, text string) (l *Lexer) {
	reader := bufio.NewReader(strings.NewReader(text))
	l, err := NewLexer(reader, "test")
	if err != nil {
		t.Fatalf("lexer instantiation failed")
	}

	return l
}

func testEq(a Token, b testToken) bool {

	if a.GetLexeme() != b.lexeme {
		return false
	}

	if a.GetTokType() != b.tokType {
		return false
	}

	if a.GetValue() != b.value {
		return false
	}

	return true
}

func TestLex(t *testing.T) {
	l := newTestLexer(t, exampleFile)
	i := 0
	for tok, err := l.Lex(); ; tok, err = l.Lex() {
		if err != nil {
			log.Fatal(err)
		}

		if tok.GetTokType() == TokEOF {
			break
		}

		if !testEq(tok, tokensFile[i]) {
			t.Errorf("TestLex failed | expected %s, got %s", tokensFile[i], tok)
		}

		i++
	}

	for i, e := range examplesList {
		l := newTestLexer(t, e)
		j := 0

		for tok, err := l.Lex(); ; tok, err = l.Lex() {
			if err != nil {
				log.Fatal(err)
			}

			if tok.GetTokType() == TokEOF {
				break
			}

			if !testEq(tok, tokensList[i][j]) {
				t.Errorf("TestLex failed | expected %s, got %s", tokensList[i][j], tok)
			}

			j++
		}
	}
}

func TestPeek(t *testing.T) {
	l := newTestLexer(t, exampleFile)

	for tok_peek, err := l.Peek(); ; tok_peek, err = l.Peek() {
		if err != nil {
			log.Fatal(err)
		}

		if tok_peek.GetTokType() == TokEOF {
			break
		}

		tok_lex, _ := l.Lex()
		if tok_peek != tok_lex {
			t.Errorf("TestPeek failed | expected %s, got %s", tok_peek, tok_lex)
		}
	}

	for _, e := range examplesList {
		l := newTestLexer(t, e)
		j := 0

		for tok_peek, err := l.Peek(); ; tok_peek, err = l.Peek() {
			if err != nil {
				log.Fatal(err)
			}

			if tok_peek.GetTokType() == TokEOF {
				break
			}

			tok_lex, _ := l.Lex()
			if tok_peek != tok_lex {
				t.Errorf("TestPeek failed | expected %s, got %s", tok_peek, tok_lex)
			}

			j++
		}
	}
}

func TestSkipUntilAndLex(t *testing.T) {
	l := newTestLexer(t, example2)
	DebugLexer = true
	l.Lex()
	l.SkipUntilAndLex(Semicolon)
	l.Lex()
	l.Lex()
	l.Lex()
}

func TestSkipUntil(t *testing.T) {
	l := newTestLexer(t, example2)
	DebugLexer = true
	l.Lex()
	l.SkipUntil(Semicolon, TokRPar)
	l.Lex()
	l.Lex()
	l.Lex()
}
