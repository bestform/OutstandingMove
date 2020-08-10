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
		t.Error("Expected 2 possible moves, but got", len(moves))
		fmt.Print(moves)
	}
}

func TestMoveOutOfChess(t *testing.T) {
	CurrentBoard, _ = fen.FenToBoard("1r2qbnr/3b1k1p/Pp3pB1/1NPp4/5NP1/4P2Q/PBP4P/R3K1R1 b Q - 0 20")

	moves := CalculatePossibleMoves(true)

	if len(moves) != 3 {
		t.Error("Expected 3 possible moves, but got", len(moves))
		fmt.Print(moves)
	}
}