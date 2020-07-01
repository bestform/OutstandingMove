package board

import (
	"strconv"
)

type File int

const (
	A File = iota
	B
	C
	D
	E
	F
	G
	H
)

var AllFiles = []File{A, B, C, D, E, F, G, H}

type Castling int

const (
	WHITE_KINGSIDE Castling = iota
	WHITE_QUEENSIDE
	BLACK_KINGSIDE
	BLACK_QUEENSIDE
)

type Position struct {
	File File
	Rank int
}

func (p Position) String() string {
	var out string
	switch p.File {
	case A:
		out += "a"
	case B:
		out += "b"
	case C:
		out += "c"
	case D:
		out += "d"
	case E:
		out += "e"
	case F:
		out += "f"
	case G:
		out += "g"
	case H:
		out += "h"
	}
	out += strconv.Itoa(p.Rank)

	return out
}

func (p *Position) SameAs(p2 *Position) bool {
	if p == nil && p2 == nil {
		return true
	}
	if p == nil || p2 == nil {
		return false
	}

	return p.File == p2.File && p.Rank == p2.Rank
}

func PosFromString(pos string) Position {
	position := Position{}
	if len(pos) != 2 {
		// we should in theory return an error, but it would make the using code very awkward,
		// so we see if we can get away with it, knowing that if there is something wrong with
		// the pos string we get really nasty behaviour.
		return position
	}
	switch pos[0] {
	case 'a':
		position.File = A
	case 'b':
		position.File = B
	case 'c':
		position.File = C
	case 'd':
		position.File = D
	case 'e':
		position.File = E
	case 'f':
		position.File = F
	case 'g':
		position.File = G
	case 'h':
		position.File = H
	}

	// same here. We really shouldn't ignore the error, but in favour of better usability we will let this one slide
	rank, _ := strconv.Atoi(string(pos[1]))
	position.Rank = rank

	return position
}

type Move struct {
	From Position
	To Position
}

func (m Move) String() string {
	var out string
	for _, p := range []Position{m.From, m.To} {
		out += p.String()
	}

	return out
}

func MoveFromString(moveStr string) Move {
	move := Move{}
	fromString := string(moveStr[0]) + string(moveStr[1])
	toString := string(moveStr[2]) + string(moveStr[3])
	move.From = PosFromString(fromString)
	move.To = PosFromString(toString)

	return move
}

type Mailbox120 [120]int
type Mailbox64 [64]int

type Cell struct {
	Occupied bool
	Occupant *Piece
}
type Board struct {
	Cells      []Cell
	Side       Color
	Castling   []Castling
	EnPassant  *Position
	HalfTurns  int
	TurnNumber int
}

func NewMailbox120() *Mailbox120 {
	return &Mailbox120{
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1,  0,  1,  2,  3,  4,  5,  6,  7, -1,
		-1,  8,  9, 10, 11, 12, 13, 14, 15, -1,
		-1, 16, 17, 18, 19, 20, 21, 22, 23, -1,
		-1, 24, 25, 26, 27, 28, 29, 30, 31, -1,
		-1, 32, 33, 34, 35, 36, 37, 38, 39, -1,
		-1, 40, 41, 42, 43, 44, 45, 46, 47, -1,
		-1, 48, 49, 50, 51, 52, 53, 54, 55, -1,
		-1, 56, 57, 58, 59, 60, 61, 62, 63, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
		-1, -1, -1, -1, -1, -1, -1, -1, -1, -1,
	}
}

func NewMailbox64() *Mailbox64 {
	return &Mailbox64{
		21, 22, 23, 24, 25, 26, 27, 28,
		31, 32, 33, 34, 35, 36, 37, 38,
		41, 42, 43, 44, 45, 46, 47, 48,
		51, 52, 53, 54, 55, 56, 57, 58,
		61, 62, 63, 64, 65, 66, 67, 68,
		71, 72, 73, 74, 75, 76, 77, 78,
		81, 82, 83, 84, 85, 86, 87, 88,
		91, 92, 93, 94, 95, 96, 97, 98,
	}
}

func NewBoard() *Board {
	return &Board{
		Cells:    make([]Cell, 64),
		Castling: make([]Castling, 0),
	}
}

func (b *Board) IsCastlingPossible(c Castling) bool {
	for _, possibleCastling := range b.Castling {
		if possibleCastling == c {
			return true
		}
	}

	return false
}

func (b *Board) SwitchSide() {
	if b.Side == WHITE {
		b.Side = BLACK
		return
	}
	b.Side = WHITE
}

func (b *Board) Move(move Move) {
	pieceToMove := b.PieceAt(move.From)
	b.ClearPieceAt(move.From)
	b.SetPieceAt(move.To, pieceToMove)
}

func (b *Board) PieceAt(p Position) *Piece {
	return b.Cells[indexFromFileAndRank(p.File, p.Rank)].Occupant
}

func (b *Board) SetPieceAt(p Position, piece *Piece) {
	b.Cells[indexFromFileAndRank(p.File, p.Rank)].Occupant = piece
	b.Cells[indexFromFileAndRank(p.File, p.Rank)].Occupied = true
}

func (b *Board) ClearPieceAt(p Position) {
	b.Cells[indexFromFileAndRank(p.File, p.Rank)].Occupant = nil
	b.Cells[indexFromFileAndRank(p.File, p.Rank)].Occupied = false
}

func indexFromFileAndRank(file File, rank int) int {
	column := int(file)

	return (rank-1)*8 + column
}

func index120FromFileAndRank(file File, rank int) int {
	column := int(file) + 1

	return (rank-1)*10 + column + 20
}

func PositionFromIndex(index int) *Position {
	p := Position{}
	p.File = File(index % 8)
	p.Rank = index / 8 + 1

	return &p
}

func (b *Board) String() string {
	str := ""
	for rank := 8; rank > 0; rank-- {
		for file := 0; file < 8; file++ {
			cell := b.PieceAt(Position{File: File(file), Rank: rank})
			if cell == nil {
				str += "."
				continue
			}
			switch cell.Kind {
			case PAWN:
				str += "P"
			case ROOK:
				str += "R"
			case BISHOP:
				str += "B"
			case QUEEN:
				str += "Q"
			case KING:
				str += "K"
			case KNIGHT:
				str += "N"
			}

		}
		str += "\n"
	}

	return str
}
