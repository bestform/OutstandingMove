package chessBot

type AN struct {
	File string
	Rank int
	Piece ChessPieceKind
	IsMove bool
	IsCapture bool
	IsPromotion bool
	PromotesTo ChessPieceKind
	IsKingsideCastling bool
	IsQueensideCastling bool
}
