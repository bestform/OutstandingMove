package board

import "fmt"

type ChessPieceKind int
type Color int

const (
	PAWN ChessPieceKind = iota
	KNIGHT
	BISHOP
	ROOK
	QUEEN
	KING
)

const (
	WHITE Color = 0
	BLACK Color = 1
)

type Piece struct {
	Kind    ChessPieceKind
	Offsets [8]int
	Slide   bool
	Color   Color
}

func (p *Piece) SameAs(p2 *Piece) bool {
	return p.Color == p2.Color && p.Kind == p2.Kind
}

func (p *Piece) String() string {
	return fmt.Sprint(p.Color, p.Kind)
}

var offsets map[ChessPieceKind][8]int
var knightAndRayDirections = [6]int{0, 8, 4, 4, 8, 8}
var canSlide = [6]bool{false, false, true, true, true, false}

func init() {
	offsets = make(map[ChessPieceKind][8]int)
	offsets[PAWN] = [8]int{}
	offsets[KNIGHT] = [8]int{-21, -19, -12, -8, 8, 12, 19, 21}
	offsets[BISHOP] = [8]int{-11, -9, 9, 11, 0, 0, 0, 0}
	offsets[ROOK] = [8]int{-10, -1, 1, 10, 0, 0, 0, 0}
	offsets[QUEEN] = [8]int{-11, -10, -9, -1, 1, 9, 10, 11}
	offsets[KING] = [8]int{-11, -10, -9, -1, 1, 9, 10, 11}
}

func NewPiece(kind ChessPieceKind, color Color) *Piece {
	p := Piece{
		Kind:    kind,
		Color:   color,
		Offsets: offsets[kind],
		Slide:   canSlide[kind],
	}

	return &p
}
