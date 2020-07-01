package engine

import (
	"chessBot/board"
	"chessBot/fen"
	"chessBot/uci"
	"fmt"
	"log"
	"os"
)

var CurrentBoard *board.Board
var LogOutput *os.File
var mb120 *board.Mailbox120
var mb64 *board.Mailbox64

func init() {
	var err error
	LogOutput, err = os.Create("outstandingMove.log")
	if err != nil {
		log.Fatal("error opening log file:", err)
	}

	mb120 = board.NewMailbox120()
	mb64 = board.NewMailbox64()
}

func Send(msg string) {
	fmt.Println(msg)
	Log(fmt.Sprintf("-> %s", msg))
}

func Log(msg string) {
	LogOutput.WriteString(msg + "\n")
	LogOutput.Sync()
}

func InitBoard(stmnt *uci.PositionStatement) error {
	initString := fen.STARTPOSSTRING

	if stmnt.IsFen {
		initString = stmnt.FenString
	}
	var err error
	CurrentBoard, err = fen.FenToBoard(initString)
	if err != nil {
		return err
	}

	for _, moveString := range stmnt.Moves {
		Log("moving " + moveString)
		move := board.MoveFromString(moveString)
		CurrentBoard.Move(move)
	}

	return nil
}

func CalculatePossibleMoves() []board.Move {
	var moves []board.Move
	posInMb64 := -1
	for rank := 1; rank < 9; rank++ {
		for _, file := range board.AllFiles {
			posInMb64++
			position := board.Position{file, rank}
			piece := CurrentBoard.PieceAt(position)
			if piece == nil {
				continue
			}
			if piece.Color != CurrentBoard.Side {
				continue
			}

			moves = append(moves, calculateMovesForPieceAt(piece, posInMb64)...)
		}
	}

	return moves
}

func calculateMovesForPieceAt(piece *board.Piece, posInMb64 int) []board.Move {
	var moves []board.Move
	from := board.PositionFromIndex(posInMb64)

	if piece.Kind == board.PAWN {
		return movesForPawn(piece, posInMb64)
	}

	for j := 0; j < piece.Directions; j++ {
		n := posInMb64
		for {
			n = mb120[mb64[n] + piece.Offsets[j]]
			if n == -1 {
				break
			}
			to := board.PositionFromIndex(n)
			targetPiece := CurrentBoard.PieceAt(*to)
			if targetPiece == nil { // todo: avoid moving into chess with the king
				moves = append(moves, board.Move{
					From: *from,
					To:   *to,
				})
				if !piece.Slide {
					break
				}
				continue
			}
			if targetPiece.Color != CurrentBoard.Side {
				moves = append(moves, board.Move{
					From: *from,
					To:   *to,
				})
			}
			break
		}
	}

	return moves
}

func movesForPawn(piece *board.Piece, posInMb64 int) []board.Move {
	var moves []board.Move
	var moveOffsets []int
	var strikeOffsets []int
	from := board.PositionFromIndex(posInMb64)
	if piece.Color == board.WHITE {
		moveOffsets = append(moveOffsets, 10)
		if from.Rank == 2 {
			moveOffsets = append(moveOffsets, 20)
		}
		strikeOffsets = []int{9, 11}
	}
	if piece.Color == board.BLACK {
		moveOffsets = append(moveOffsets, -10)
		if from.Rank == 7 {
			moveOffsets = append(moveOffsets, -20)
		}
		strikeOffsets = []int{-9, -11}
	}

	// moves
	for _, offset := range moveOffsets {
		newPosInt := mb120[mb64[posInMb64] + offset]
		if newPosInt == -1 {
			continue // todo: Handle promotion
		}
		to := board.PositionFromIndex(newPosInt)
		pieceAtNewPos := CurrentBoard.PieceAt(*to)
		if pieceAtNewPos != nil {
			break
		}
		moves = append(moves, board.Move{
			From: *from,
			To:   *to,
		})
	}

	// strikes
	for _, offset := range strikeOffsets {
		newPosInt := mb120[mb64[posInMb64] + offset]
		if newPosInt == -1 {
			continue
		}
		to := board.PositionFromIndex(newPosInt)
		pieceAtNewPos := CurrentBoard.PieceAt(*to)
		if pieceAtNewPos == nil {
			continue
		}
		if pieceAtNewPos.Color != CurrentBoard.Side {
			moves = append(moves, board.Move{
				From: *from,
				To:   *to,
			})
		}
	}

	return moves
}
