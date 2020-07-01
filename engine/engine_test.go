package engine

import (
	"chessBot/fen"
	"fmt"
	"testing"
)

func TestCalculatePossibleMoves(t *testing.T) {
	CurrentBoard, _ = fen.FenToBoard(fen.STARTPOSFEN)

	moves := CalculatePossibleMoves()
	fmt.Println(len(moves))
}