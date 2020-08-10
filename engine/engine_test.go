package engine

import (
	"chessBot/board"
	"chessBot/fen"
	"fmt"
	"testing"
)

func TestCalculatePossibleMoves(t *testing.T) {
	CurrentBoard, _ = fen.FenToBoard(fen.STARTPOSFEN)
	moves := CalculatePossibleMoves(true)

	if len(moves) != 20 {
		t.Error("Expected 20 possible moves, but got", len(moves))
		fmt.Print(moves)
	}
}

func TestCalculatePossibleMovesForBlack(t *testing.T) {
	CurrentBoard, _ = fen.FenToBoard(fen.STARTPOSFEN)
	CurrentBoard.Side = board.BLACK
	moves := CalculatePossibleMoves(true)

	if len(moves) != 20 {
		t.Error("Expected 20 possible moves, but got", len(moves))
		fmt.Print(moves)
	}
}

func TestCalculatePossibleMovesNoMoveToChess(t *testing.T) {
	CurrentBoard, _ = fen.FenToBoard("4k3/8/7b/7b/8/8/4P3/3K4 w - - 0 1")

	moves := CalculatePossibleMoves(true)

	if len(moves) != 2 {
		t.Error("Expected 1 possible move, but got", len(moves))
		fmt.Print(moves)
	}

}