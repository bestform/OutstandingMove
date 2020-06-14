package uci

import (
	"fmt"
	"strings"
)

type StatementKind int

const (
	UciStatementKind StatementKind = iota
	DebugStatementKind
	IsReadyStatementKind
	SetOptionStatementKind
	RegisterStatementKind
	UciNewGameStatementKind
	PositionStatementKind
	GoStatementKind
	Go_searchMovesStatementKind
	Go_ponderStatementKind
	Go_wtimeStatementKind
	Go_btimeStatementKind
	Go_wincStatementKind
	Go_bincStatementKind
	Go_movesToGoStatementKind
	Go_depthStatementKind
	Go_nodesStatementKind
	Go_mateStatementKind
	Go_moveTimeStatementKind
	Go_inifiniteStatementKind
	StopStatementKind
	PonderHitStatementKind
	QuitStatementKind
)

type Statement struct {
	Kind      StatementKind
	Debug     *DebugStatement
	SetOption *SetOptionStatement
	Register  *RegisterStatement
	Position  *PositionStatement
	Go        *GoStatement
}

type SetOptionStatement struct{
	Name string
	Value string
}

type RegisterStatement struct{}

type PositionStatement struct{}

type GoStatement struct{}

type DebugStatement struct{
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
