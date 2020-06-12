package uci

import (
	"testing"
)

func TestAlgebraicNotation(t *testing.T) {
	tokens, _ := Lex("e8f8 a7e2q")

	expectedTokens := []*token{
		{
			kind:  longAlgebraicNotation,
			value: "e8f8",
		},
		{
			kind:  longAlgebraicNotation,
			value: "a7e2q",
		},
	}
	for i, actual := range tokens {
		if !actual.equals(expectedTokens[i]) {
			t.Errorf("Expected %v, but got %v", expectedTokens[i], actual)
		}
	}
}

func TestKeywords(t *testing.T) {
	keywords := []string{"value", "binc", "btime", "code", "debug", "depth", "fen", "go", "infinite", "isready", "later", "mate", "moves", "movestogo", "movetime", "name", "nodes", "off", "on", "ponder", "ponderhit", "position", "quit", "register", "searchmoves", "setoption", "startpos", "stop", "uci", "ucinewgame", "winc", "wtime"}

	var input string
	for _, keywrd := range keywords {
		input += " " + keywrd
	}

	actualList, _ := Lex(input)

	if len(actualList) != len(keywords) {
		t.Error("Unexpected number of tokens")
	}

	for i, actualToken := range actualList {
		if actualToken.kind != keywordKind {
			t.Errorf("Expected %v to have kind %d, but got %d", actualToken, keywordKind, actualToken.kind)
		}
		if actualToken.value != keywords[i] {
			t.Errorf("Expected value to be %s but got %s", keywords[i], actualToken.value)
		}
	}
}

func TestIgnoreSpace(t *testing.T) {
	actual, _ := Lex("go          e8f8 e8f8")

	if len(actual) != 3 {
		t.Errorf("Expected 3 tokens but got %d", len(actual))
	}
}

func TestSymbols(t *testing.T) {
	actualList, _ := Lex("go;e8f8\n")

	if len(actualList) != 4 {
		t.Error("Expected 4 tokens, but got", len(actualList))
	}

	expectedSemicolon := &token{symbolKind, ";"}
	if !actualList[1].equals(expectedSemicolon) {
		t.Errorf("Expected %v, but got %v", expectedSemicolon, actualList[1])
	}

	expectedNewLine := &token{symbolKind, "\n"}
	if !actualList[3].equals(expectedNewLine) {
		t.Errorf("Expected %v, but got %v", expectedNewLine, actualList[3])
	}
}

func TestStrings(t *testing.T) {
	actualList, _ := Lex("string1 c:\\path\\string2;string3")

	if len(actualList) != 4 {
		t.Error("Expected 4 tokens, but got", len(actualList))
	}

	for _, i := range []int{0,1,3} {
		if actualList[i].kind != stringKind {
			t.Errorf("Expected %v to have kind %d, but got %d", actualList[i], stringKind, actualList[i].kind)
		}
	}

	if actualList[0].value != "string1" {
		t.Error("Expected value string1, but got", actualList[0].value)
	}
	if actualList[0].value != "string1" {
		t.Error("Expected value string1, but got", actualList[0].value)
	}
	if actualList[1].value != "c:\\path\\string2" {
		t.Error("Expected value c:\\path\\string2, but got", actualList[1].value)
	}
	if actualList[3].value != "string3" {
		t.Error("Expected value string3, but got", actualList[3].value)
	}
}


