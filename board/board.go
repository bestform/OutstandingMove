package board

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

type Mailbox120 [120]int
type Cell struct {
	Occupied bool
	Occupant *Piece
}
type Board struct {
	Cells []Cell
}


func NewMailbox120() *Mailbox120 {
	return &Mailbox120 {
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

func NewBoard() *Board {
	return &Board{
		Cells: make([]Cell, 64),
	}
}

func (b *Board) CellAt(file File, rank int) *Piece {
	return b.Cells[positionFromFileAndRank(file, rank)].Occupant
}

func (b *Board) SetCellAt(file File, rank int, piece *Piece) {
	b.Cells[positionFromFileAndRank(file, rank)].Occupant = piece
	b.Cells[positionFromFileAndRank(file, rank)].Occupied = true
}

func (b *Board) ClearCellAt(file File, rank int) {
	b.Cells[positionFromFileAndRank(file, rank)].Occupant = nil
	b.Cells[positionFromFileAndRank(file, rank)].Occupied = false
}

func positionFromFileAndRank(file File, rank int) int {
	column := int(file)

	return (rank - 1) * 8 + column
}

func position120FromFileAndRank(file File, rank int) int {
	column := int(file) + 1

	return (rank - 1) * 10 + column + 20
}

func (b *Board) String() string {
	str := ""
	for rank := 8; rank > 0; rank-- {
		for file := 0; file < 8; file++ {
			cell := b.CellAt(File(file), rank)
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




