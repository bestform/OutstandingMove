package uci

import (
	"fmt"
	"strconv"
	"strings"
)

type StatementKind string

const (
	UciStatementKind StatementKind = "uci"
	DebugStatementKind = "debug"
	IsReadyStatementKind = "isReady"
	SetOptionStatementKind = "setOption"
	RegisterStatementKind = "register"
	UciNewGameStatementKind = "uciNewGame"
	PositionStatementKind = "position"
	GoStatementKind = "go"
	StopStatementKind = "stop"
	PonderHitStatementKind = "ponderHit"
	QuitStatementKind = "quit"
)

type GoKind int

const (
	Go_searchMovesKind GoKind = iota
	Go_ponderKind
	Go_wtimeKind
	Go_btimeKind
	Go_wincKind
	Go_bincKind
	Go_movesToGoKind
	Go_depthKind
	Go_nodesKind
	Go_mateKind
	Go_moveTimeKind
	Go_inifiniteKind
)

type Statement struct {
	Kind      StatementKind
	Debug     *DebugStatement
	SetOption *SetOptionStatement
	Register  *RegisterStatement
	Position  *PositionStatement
	Go        *GoStatement
}

type SetOptionStatement struct {
	Name  string
	Value string
}

type RegisterStatement struct {
	IsLater bool
	Name    string
	Code    string
}

type PositionStatement struct {
	IsFen      bool
	IsStartPos bool
	FenString  string
	Moves      []string
}

type GoStatement struct {
	Kinds       []GoKind
	SearchMoves []string
	Wtime       int
	Btime       int
	Winc        int
	Binc        int
	MovesToGo   int
	Depth       int
	Nodes       int
	Mate        int
	MoveTime    int
}

type DebugStatement struct {
	On bool
}

func Parse(source string) ([]*Statement, error) {
	var statements []*Statement
	tokens, err := Lex(source)
	if err != nil {
		return statements, err
	}

	var cursor uint

	for cursor < uint(len(tokens)) {
		cursor = eatNonKeywords(tokens, cursor)
		if cursor == uint(len(tokens))-1 {
			break
		}

		// uci statement
		newCursor, ok := parseSingleKeywordStatement(tokens, cursor, uci)
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind: UciStatementKind,
			})
			if newCursor == uint(len(tokens)) {
				break
			}
		}

		// debug statement
		debugStmnt, newCursor, ok, err := parseDebugStatement(tokens, cursor)
		if err != nil {
			return statements, fmt.Errorf("error parsing debug statement: %s", err)
		}
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind:  DebugStatementKind,
				Debug: debugStmnt,
			})
			if newCursor == uint(len(tokens)) {
				break
			}
		}

		//isReady statement
		newCursor, ok = parseSingleKeywordStatement(tokens, cursor, isready)
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind: IsReadyStatementKind,
			})
			if newCursor == uint(len(tokens)) {
				break
			}
		}

		// setOption statement
		setOptionStmnt, newCursor, ok, err := parseSetOptionStatement(tokens, cursor)
		if err != nil {
			return statements, fmt.Errorf("error parsing setOption: %s", err)
		}
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind:      SetOptionStatementKind,
				SetOption: setOptionStmnt,
			})
			if newCursor == uint(len(tokens)) {
				break
			}
		}

		// position statement
		positionStmnt, newCursor, ok, err := parsePositionStatement(tokens, cursor)
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind:     PositionStatementKind,
				Position: positionStmnt,
			})
			if newCursor == uint(len(tokens)) {
				break
			}
		}

		// go statement
		goStmnt, newCursor, ok, err := parseGoStatement(tokens, cursor)
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind: GoStatementKind,
				Go:   goStmnt,
			})
			if newCursor == uint(len(tokens)) {
				break
			}
		}

		// register statement
		registerStmnt, newCursor, ok, err := parseRegisterStatement(tokens, cursor)
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind:     RegisterStatementKind,
				Register: registerStmnt,
			})
			if newCursor == uint(len(tokens)) {
				break
			}
		}

		// newGame Statement
		newCursor, ok = parseSingleKeywordStatement(tokens, cursor, ucinewgame)
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind: UciNewGameStatementKind,
			})
			if newCursor == uint(len(tokens)) {
				break
			}
		}

		// stop Statement
		newCursor, ok = parseSingleKeywordStatement(tokens, cursor, stop)
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind: StopStatementKind,
			})
			if newCursor == uint(len(tokens)) {
				break
			}
		}

		// ponderHit Statement
		newCursor, ok = parseSingleKeywordStatement(tokens, cursor, ponderhit)
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind: PonderHitStatementKind,
			})
			if newCursor == uint(len(tokens)) {
				break
			}
		}

		// quit Statement
		newCursor, ok = parseSingleKeywordStatement(tokens, cursor, quit)
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind: QuitStatementKind,
			})
			if newCursor == uint(len(tokens)) {
				break
			}
		}
	}

	return statements, nil
}

func parseRegisterStatement(tokens []*token, ic uint) (*RegisterStatement, uint, bool, error) {
	if !tokens[ic].equals(tokenFromKeyword(go_)) {
		return nil, ic, false, nil
	}
	cursor := ic

	stmnt := &RegisterStatement{}
	cursor++

	if tokens[cursor].equals(tokenFromKeyword(later)) {
		stmnt.IsLater = true
		cursor = eatAllToNextNewLine(tokens, cursor)
		return stmnt, cursor, true, nil
	}

	for !tokens[cursor].equals(tokenFromSymbol(newLine)) && cursor < uint(len(tokens)) {
		if tokens[cursor].equals(tokenFromKeyword(name)) {
			var name []string
			for !tokens[cursor].equals(tokenFromKeyword(code)) && !tokens[cursor].equals(tokenFromSymbol(newLine)) && cursor < uint(len(tokens)) {
				name = append(name, tokens[cursor].value)
				cursor++
			}
			stmnt.Name = strings.Join(name, " ")
		}
		if tokens[cursor].equals(tokenFromKeyword(code)) {
			var code []string
			for !tokens[cursor].equals(tokenFromKeyword(name)) && !tokens[cursor].equals(tokenFromSymbol(newLine)) && cursor < uint(len(tokens)) {
				code = append(code, tokens[cursor].value)
				cursor++
			}
			stmnt.Code = strings.Join(code, " ")
		}
	}

	return stmnt, cursor, true, nil
}

func parseGoStatement(tokens []*token, ic uint) (*GoStatement, uint, bool, error) {
	if !tokens[ic].equals(tokenFromKeyword(go_)) {
		return nil, ic, false, nil
	}
	cursor := ic

	stmnt := &GoStatement{
		Kinds: make([]GoKind, 0),
	}

	parseIntGoStatement := func(kind GoKind, stmnt *GoStatement, ic uint) (int, uint, error) {
		stmnt.Kinds = append(stmnt.Kinds, kind)
		cursor := ic
		cursor++
		number, err := strconv.Atoi(tokens[cursor].value)
		if err != nil {
			return 0, ic, fmt.Errorf("expected number: %s", err)
		}
		cursor++
		return number, cursor, nil
	}

	var err error
	var number int

	cursor++
	for cursor < uint(len(tokens)) && !tokens[cursor].equals(tokenFromSymbol(newLine)) {

		switch tokens[cursor].value {

		case string(searchmoves):
			cursor++
			stmnt.Kinds = append(stmnt.Kinds, Go_searchMovesKind)
			var moves []string
			for tokens[cursor].kind == longAlgebraicNotation {
				moves = append(moves, tokens[cursor].value)
				cursor++
			}
			stmnt.SearchMoves = moves

		case string(ponder):
			stmnt.Kinds = append(stmnt.Kinds, Go_ponderKind)
			cursor++

		case string(infinite):
			stmnt.Kinds = append(stmnt.Kinds, Go_inifiniteKind)
			cursor++

		case string(wtime):
			number, cursor, err = parseIntGoStatement(Go_wtimeKind, stmnt, cursor)
			if err != nil {
				return nil, ic, false, err
			}
			stmnt.Wtime = number

		case string(btime):
			number, cursor, err = parseIntGoStatement(Go_btimeKind, stmnt, cursor)
			if err != nil {
				return nil, ic, false, err
			}
			stmnt.Btime = number

		case string(winc):
			number, cursor, err = parseIntGoStatement(Go_wincKind, stmnt, cursor)
			if err != nil {
				return nil, ic, false, err
			}
			stmnt.Winc = number

		case string(binc):
			number, cursor, err = parseIntGoStatement(Go_bincKind, stmnt, cursor)
			if err != nil {
				return nil, ic, false, err
			}
			stmnt.Binc = number

		case string(movestogo):
			number, cursor, err = parseIntGoStatement(Go_movesToGoKind, stmnt, cursor)
			if err != nil {
				return nil, ic, false, err
			}
			stmnt.MovesToGo = number

		case string(depth):
			number, cursor, err = parseIntGoStatement(Go_depthKind, stmnt, cursor)
			if err != nil {
				return nil, ic, false, err
			}
			stmnt.Depth = number

		case string(nodes):
			number, cursor, err = parseIntGoStatement(Go_nodesKind, stmnt, cursor)
			if err != nil {
				return nil, ic, false, err
			}
			stmnt.Nodes = number

		case string(mate):
			number, cursor, err = parseIntGoStatement(Go_mateKind, stmnt, cursor)
			if err != nil {
				return nil, ic, false, err
			}
			stmnt.Mate = number

		case string(movetime):
			number, cursor, err = parseIntGoStatement(Go_moveTimeKind, stmnt, cursor)
			if err != nil {
				return nil, ic, false, err
			}
			stmnt.MoveTime = number
		}
	}
	cursor++

	return stmnt, cursor, true, nil
}

func parsePositionStatement(tokens []*token, ic uint) (*PositionStatement, uint, bool, error) {
	if !tokens[ic].equals(tokenFromKeyword(position)) {
		return nil, ic, false, nil
	}
	cursor := ic
	cursor++
	stmnt := &PositionStatement{}
	if tokens[cursor].equals(tokenFromKeyword(fen)) {
		stmnt.IsFen = true
		cursor++
		stmnt.FenString = tokens[cursor].value
		cursor++
	}
	if tokens[cursor].equals(tokenFromKeyword(startpos)) {
		stmnt.IsStartPos = true
		if stmnt.IsFen {
			return nil, ic, false, fmt.Errorf("a position statement cannot have a fen string and a startpos flag")
		}
		cursor++
	}

	if !tokens[cursor].equals(tokenFromKeyword(moves)) {
		return nil, ic, false, fmt.Errorf("expected 'moves' in position statement")
	}
	cursor++
	var moves []string
	for !tokens[cursor].equals(tokenFromSymbol(newLine)) && cursor < uint(len(tokens)) {
		if tokens[cursor].kind != longAlgebraicNotation {
			return nil, ic, false, fmt.Errorf("expected long algebraic notation string for moves, but found %s", tokens[cursor].value)
		}
		moves = append(moves, tokens[cursor].value)
		cursor++
	}
	stmnt.Moves = moves
	cursor++

	return stmnt, cursor, true, nil
}

func parseSetOptionStatement(tokens []*token, cursor uint) (*SetOptionStatement, uint, bool, error) {
	if !tokens[cursor].equals(tokenFromKeyword(setoption)) {
		return nil, cursor, false, nil
	}

	cursor++
	if !tokens[cursor].equals(tokenFromKeyword(name)) {
		return nil, cursor, false, fmt.Errorf("expected name after setOption")
	}

	stmnt := &SetOptionStatement{}

	cursor++
	var nameParts []string
	for !tokens[cursor].equals(tokenFromKeyword(value)) && !tokens[cursor].equals(tokenFromSymbol(newLine)) {
		nameParts = append(nameParts, tokens[cursor].value)
		cursor++
	}
	stmnt.Name = strings.Join(nameParts, " ")

	if !tokens[cursor].equals(tokenFromKeyword(value)) {
		cursor++
		return stmnt, cursor, true, nil
	}
	cursor++
	var valueParts []string
	for !tokens[cursor].equals(tokenFromSymbol(newLine)) {
		valueParts = append(valueParts, tokens[cursor].value)
		cursor++
	}
	stmnt.Value = strings.Join(valueParts, " ")
	cursor++

	return stmnt, cursor, true, nil
}

func parseDebugStatement(tokens []*token, cursor uint) (*DebugStatement, uint, bool, error) {
	if !tokens[cursor].equals(tokenFromKeyword(debug)) {
		return nil, cursor, false, nil
	}

	cursor++
	if !tokens[cursor].equals(tokenFromKeyword(on)) || !tokens[cursor].equals(tokenFromKeyword(off)) {
		return nil, cursor, false, fmt.Errorf("expected on or off after debug statement")
	}
	onOffToken := tokens[cursor]

	cursor = eatAllToNextNewLine(tokens, cursor)

	return &DebugStatement{On: onOffToken.equals(tokenFromKeyword(on))}, cursor, true, nil
}

func tokenFromKeyword(k keyword) *token {
	return &token{
		kind:  keywordKind,
		value: string(k),
	}
}

func tokenFromSymbol(s symbol) *token {
	return &token{
		kind:  symbolKind,
		value: string(s),
	}
}

func eatAllToNextNewLine(tokens []*token, cursor uint) uint {
	for cursor < uint(len(tokens))-1 && !tokens[cursor].equals(tokenFromSymbol(newLine)) {
		cursor++
	}

	return cursor + 1
}

func parseSingleKeywordStatement(tokens []*token, cursor uint, keyword keyword) (uint, bool) {
	if !tokens[cursor].equals(tokenFromKeyword(keyword)) {
		return cursor, false
	}

	cursor++
	cursor = eatAllToNextNewLine(tokens, cursor)

	return cursor, true
}

func eatNonKeywords(tokens []*token, cursor uint) uint {
	for uint(len(tokens)) < cursor && tokens[cursor].kind != keywordKind {
		cursor++
	}

	return cursor
}
