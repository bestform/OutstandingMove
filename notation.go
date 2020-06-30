package chessBot

import "chessBot/board"

type AN struct {
	File string
	Rank int
	Piece board.ChessPieceKind
	IsMove bool
	IsCapture bool
	IsPromotion bool
	PromotesTo board.ChessPieceKind
	IsKingsideCastling bool
	IsQueensideCastling bool
}
