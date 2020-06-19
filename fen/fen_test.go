package fen

import (
	"fmt"
	"testing"
)

func TestFenToBoard(t *testing.T) {
	board, err := FenToBoard("rnbqkbnr/pppppppp/2p5/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(board)
}
