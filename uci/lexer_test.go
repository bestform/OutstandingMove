package uci

import (
	"fmt"
	"testing"
)

func TestLexAlgebraicNotation(t *testing.T) {
	tokens, _ := Lex("e8f8 a7e2q go foobar\n")

	for _, t := range tokens {
		fmt.Println(t)
	}

}