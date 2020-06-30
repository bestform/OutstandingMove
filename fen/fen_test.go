package fen

import (
	"chessBot/board"
	"testing"
)

func TestFenToBoard(t *testing.T) {
	createdBoard, err := FenToBoard("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		t.Fatal(err)
	}

	assertStartPosition(createdBoard, t)
}

func TestFenSpecialStringToBoard(t *testing.T) {
	createdBoard, err := FenToBoard("startpos")
	if err != nil {
		t.Fatal(err)
	}

	assertStartPosition(createdBoard, t)
}

func TestCastling(t *testing.T) {

	for _, testCase := range []struct{Fen string; Expected []board.Castling}{
		{
			Fen:      "8/8/8/8/8/8/8/8 w KQkq - 0 1",
			Expected: []board.Castling{board.WHITE_KINGSIDE, board.WHITE_QUEENSIDE, board.BLACK_KINGSIDE, board.BLACK_QUEENSIDE},
		},
		{
			Fen:      "8/8/8/8/8/8/8/8 w Kq - 0 1",
			Expected: []board.Castling{board.WHITE_KINGSIDE, board.BLACK_QUEENSIDE},
		},
		{
			Fen:      "8/8/8/8/8/8/8/8 w - - 0 1",
			Expected: []board.Castling{},
		},
	} {
		createdBoard, err := FenToBoard(testCase.Fen)
		if err != nil {
			t.Fatal(err)
		}

		if len(createdBoard.Castling) != len(testCase.Expected) {
			t.Error("Expected", len(testCase.Expected), "options, but got", len(createdBoard.Castling))
		}

		for _, expectedCastling := range testCase.Expected {
			if !createdBoard.IsCastlingPossible(expectedCastling) {
				t.Error("Expected", expectedCastling, "to be possible, but it is not")
			}
		}
	}
}

func TestHalfTurns(t *testing.T) {
	createdBoard, err := FenToBoard("8/8/8/8/8/8/8/8 w KQkq - 10 1")
	if err != nil {
		t.Fatal(err)
	}

	if createdBoard.HalfTurns != 10 {
		t.Error("Expected 10 HalfTurns but got", createdBoard.HalfTurns)
	}
}

func TestTurnNumber(t *testing.T) {
	createdBoard, err := FenToBoard("8/8/8/8/8/8/8/8 w KQkq - 1 10")
	if err != nil {
		t.Fatal(err)
	}

	if createdBoard.TurnNumber != 10 {
		t.Error("Expected TurnNumber 10 but got", createdBoard.TurnNumber)
	}
}

func TestEnPasset(t *testing.T) {

	for _, testCase := range []struct{Fen string; Expected *board.Position}{
		{
			Fen:      "8/8/8/8/8/8/8/8 w KQkq - 0 1",
			Expected: nil,
		},
		{
			Fen:      "8/8/8/8/8/8/8/8 w KQkq f3 0 1",
			Expected: &board.Position{
				File: board.F,
				Rank: 3,
			},
		},
	} {
		createdBoard, err := FenToBoard(testCase.Fen)
		if err != nil {
			t.Fatal(err)
		}

		if testCase.Expected == nil && createdBoard.EnPassant != nil {
			t.Error("Expected no EnPasset but got", createdBoard.EnPassant)
		}

		if !testCase.Expected.SameAs(createdBoard.EnPassant) {
			t.Error("Unexpected EnPasset")
		}
	}
}

func TestTurn(t *testing.T) {
	createdBoard, err := FenToBoard("8/8/8/8/8/8/8/8 w KQkq - 0 1")
	if err != nil {
		t.Fatal(err)
	}

	if createdBoard.Turn != board.WHITE {
		t.Error("Expected white's turn")
	}

	createdBoard, err = FenToBoard("8/8/8/8/8/8/8/8 b KQkq - 0 1")
	if err != nil {
		t.Fatal(err)
	}

	if createdBoard.Turn != board.BLACK {
		t.Error("Expected black's turn")
	}
}

func assertStartPosition(boardToTest *board.Board, t *testing.T) {
	for f, p := range map[board.File]*board.Piece{
		board.A: {Kind: board.ROOK, Color: board.WHITE},
		board.B: {Kind: board.KNIGHT, Color: board.WHITE},
		board.C: {Kind: board.BISHOP, Color: board.WHITE},
		board.D: {Kind: board.QUEEN, Color: board.WHITE},
		board.E: {Kind: board.KING, Color: board.WHITE},
		board.F: {Kind: board.BISHOP, Color: board.WHITE},
		board.G: {Kind: board.KNIGHT, Color: board.WHITE},
		board.H: {Kind: board.ROOK, Color: board.WHITE},
	} {
		if !boardToTest.PieceAt(board.Position{File: f, Rank: 1}).SameAs(p) {
			t.Error("expected", p, "at position", f, 1, "but got", boardToTest.PieceAt(board.Position{File: f, Rank: 1}))
		}
	}
	for f, p := range map[board.File]*board.Piece{
		board.A: {Kind: board.PAWN, Color: board.WHITE},
		board.B: {Kind: board.PAWN, Color: board.WHITE},
		board.C: {Kind: board.PAWN, Color: board.WHITE},
		board.D: {Kind: board.PAWN, Color: board.WHITE},
		board.E: {Kind: board.PAWN, Color: board.WHITE},
		board.F: {Kind: board.PAWN, Color: board.WHITE},
		board.G: {Kind: board.PAWN, Color: board.WHITE},
		board.H: {Kind: board.PAWN, Color: board.WHITE},
	} {
		if !boardToTest.PieceAt(board.Position{File: f, Rank: 2}).SameAs(p) {
			t.Error("expected", p, "at position", f, 2, "but got", boardToTest.PieceAt(board.Position{File: f, Rank: 2}))
		}
	}
	for f, p := range map[board.File]*board.Piece{
		board.A: {Kind: board.ROOK, Color: board.BLACK},
		board.B: {Kind: board.KNIGHT, Color: board.BLACK},
		board.C: {Kind: board.BISHOP, Color: board.BLACK},
		board.D: {Kind: board.QUEEN, Color: board.BLACK},
		board.E: {Kind: board.KING, Color: board.BLACK},
		board.F: {Kind: board.BISHOP, Color: board.BLACK},
		board.G: {Kind: board.KNIGHT, Color: board.BLACK},
		board.H: {Kind: board.ROOK, Color: board.BLACK},
	} {
		if !boardToTest.PieceAt(board.Position{File: f, Rank: 8}).SameAs(p) {
			t.Error("expected", p, "at position", f, 8, "but got", boardToTest.PieceAt(board.Position{File: f, Rank: 8}))
		}
	}
	for f, p := range map[board.File]*board.Piece{
		board.A: {Kind: board.PAWN, Color: board.BLACK},
		board.B: {Kind: board.PAWN, Color: board.BLACK},
		board.C: {Kind: board.PAWN, Color: board.BLACK},
		board.D: {Kind: board.PAWN, Color: board.BLACK},
		board.E: {Kind: board.PAWN, Color: board.BLACK},
		board.F: {Kind: board.PAWN, Color: board.BLACK},
		board.G: {Kind: board.PAWN, Color: board.BLACK},
		board.H: {Kind: board.PAWN, Color: board.BLACK},
	} {
		if !boardToTest.PieceAt(board.Position{File: f, Rank: 7}).SameAs(p) {
			t.Error("expected", p, "at position", f, 8, "but got", boardToTest.PieceAt(board.Position{File: f, Rank: 7}))
		}
	}

	for i := 3; i < 7; i++ {
		for _, f := range []board.File{board.A, board.B, board.C, board.D, board.E, board.F, board.G} {
			if boardToTest.PieceAt(board.Position{File: f, Rank: i}) != nil {
				t.Error("expected empty cell at", f, i, "but got", boardToTest.PieceAt(board.Position{File: f, Rank: i}))
			}
		}
	}
}
