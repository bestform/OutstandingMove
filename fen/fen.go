package fen

import (
	"chessBot/board"
	"strconv"
	"strings"
)


func FenToBoard(fenString string) (*board.Board, error) {
	parts := strings.Split(fenString, " ")
	piecesPart := parts[0]

	ranks := strings.Split(piecesPart, "/")
	newBoard := board.NewBoard()

	for oneLessCurrentRank, row := range ranks {
		var file int32
		for _, c := range row {
			if file > 7 {
				break
			}
			if c > '0' && c < '9' {
				toAdd, err := strconv.Atoi(string(c))
				if err != nil {
					return nil, err
				}
				file += int32(toAdd)
				continue
			}

			var kind board.ChessPieceKind
			isLower := string(c) == strings.ToLower(string(c))
			var color board.Color
			if isLower {
				color = board.WHITE
			} else {
				color = board.BLACK
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
			newBoard.SetCellAt(board.File(file), oneLessCurrentRank + 1, board.NewPiece(kind, color))
			file++
		}
	}

	return newBoard, nil
}

