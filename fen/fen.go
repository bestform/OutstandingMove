package fen

import (
	"chessBot/board"
	"errors"
	"strconv"
	"strings"
)

const STARTPOSSTRING = "startpos"
const STARTPOSFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func FenToBoard(fenString string) (*board.Board, error) {
	if fenString == STARTPOSSTRING {
		fenString = STARTPOSFEN
	}

	newBoard := board.NewBoard()

	parts := strings.Split(fenString, " ")
	piecesPart := parts[0]

	err := addPieces(newBoard, piecesPart)
	if err != nil {
		return nil, err
	}

	turnPart := parts[1]
	err = addTurn(newBoard, turnPart)
	if err != nil {
		return nil, err
	}

	castlingPart := parts[2]
	err = addCastling(newBoard, castlingPart)
	if err != nil {
		return nil, err
	}

	enPassetPart := parts[3]
	err = addEnPasset(newBoard, enPassetPart)
	if err != nil {
		return nil, err
	}

	halfTurnsPart := parts[4]
	err = addHalfTurns(newBoard, halfTurnsPart)
	if err != nil {
		return nil, err
	}

	turnNumberPart := parts[5]
	err = addTurnNumber(newBoard, turnNumberPart)
	if err != nil {
		return nil, err
	}

	return newBoard, nil
}

func addHalfTurns(newBoard *board.Board, part string) error {
	t, err := strconv.Atoi(part)
	if err != nil {
		return err
	}
	newBoard.HalfTurns = t

	return nil
}

func addTurnNumber(newBoard *board.Board, part string) error {
	t, err := strconv.Atoi(part)
	if err != nil {
		return err
	}
	newBoard.TurnNumber = t

	return nil
}

func addEnPasset(newBoard *board.Board, enPassetPart string) error {
	if enPassetPart == "-" {
		return nil
	}
	pos := board.PosFromString(enPassetPart)
	newBoard.EnPassant = &pos

	return nil
}

func addCastling(newBoard *board.Board, castlingPart string) error {
	if castlingPart == "-" {
		return nil
	}

	for _, c := range castlingPart {
		switch c {
		case 'K':
			newBoard.Castling = append(newBoard.Castling, board.WHITE_KINGSIDE)
		case 'Q':
			newBoard.Castling = append(newBoard.Castling, board.WHITE_QUEENSIDE)
		case 'k':
			newBoard.Castling = append(newBoard.Castling, board.BLACK_KINGSIDE)
		case 'q':
			newBoard.Castling = append(newBoard.Castling, board.BLACK_QUEENSIDE)
		default:
			return errors.New("Invalid castling char: " + string(c))
		}
	}

	return nil
}

func addTurn(newBoard *board.Board, turnPart string) error {
	switch turnPart {
	case "w":
		newBoard.Side = board.WHITE
	case "b":
		newBoard.Side = board.BLACK
	default:
		return errors.New("invalid turn part")
	}

	return nil
}

func addPieces(newBoard *board.Board, piecesPart string) error{
	ranks := strings.Split(piecesPart, "/")

	for i, row := range ranks {
		rank := 8 - i
		var file int32
		for _, c := range row {
			if file > 7 {
				break
			}
			if c > '0' && c < '9' {
				toAdd, err := strconv.Atoi(string(c))
				if err != nil {
					return err
				}
				file += int32(toAdd)
				continue
			}

			var kind board.ChessPieceKind
			isLower := string(c) == strings.ToLower(string(c))
			var color board.Color
			if isLower {
				color = board.BLACK
			} else {
				color = board.WHITE
			}

			switch strings.ToLower(string(c)) {
			case "r":
				kind = board.ROOK
			case "n":
				kind = board.KNIGHT
			case "b":
				kind = board.BISHOP
			case "q":
				kind = board.QUEEN
			case "k":
				kind = board.KING
			case "p":
				kind = board.PAWN
			}
			newBoard.SetPieceAt(board.Position{File: board.File(file), Rank: rank}, board.NewPiece(kind, color))
			file++
		}
	}

	return nil
}

