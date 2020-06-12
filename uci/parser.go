package uci

import "fmt"

type StatementKind int

const (
	uciStatement StatementKind = iota
	debugStatement
	isReadyStatement
	setOptionStatement
	registerStatement
	uciNewGameStatement
	positionStatement
	goStatement
	go_searchMovesStatement
	go_ponderStatement
	go_wtimeStatement
	go_btimeStatement
	go_wincStatement
	go_bincStatement
	go_movesToGoStatement
	go_depthStatement
	go_nodesStatement
	go_mateStatement
	go_moveTimeStatement
	go_inifiniteStatement
	stopStatement
	ponderHitStatement
	quitStatement
)

type Statement struct {
	Kind      StatementKind
	Debug     *DebugStatement
	SetOption *SetOptionStatement
	Register  *RegisterStatement
	Position  *PositionStatement
	Go        *GoStatement
}

type SetOptionStatement struct{}

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
				Kind: uciStatement,
			})
		}

		// debug statement
		debugStmnt, newCursor, ok, err := parseDebugStatement(tokens, cursor)
		if err != nil {
			return statements, fmt.Errorf("error parsing debug statement: %s", err)
		}
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind:  debugStatement,
				Debug: debugStmnt,
			})
		}

		//isReady statement
		newCursor, ok = parseSingleKeywordStatement(tokens, cursor, isready)
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind: isReadyStatement,
			})
		}

		// newGame Statement
		newCursor, ok = parseSingleKeywordStatement(tokens, cursor, ucinewgame)
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind: uciNewGameStatement,
			})
		}

		// stop Statement
		newCursor, ok = parseSingleKeywordStatement(tokens, cursor, stop)
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind: stopStatement,
			})
		}

		// ponderHit Statement
		newCursor, ok = parseSingleKeywordStatement(tokens, cursor, ponderhit)
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind: ponderHitStatement,
			})
		}

		// quit Statement
		newCursor, ok = parseSingleKeywordStatement(tokens, cursor, quit)
		if ok {
			cursor = newCursor
			statements = append(statements, &Statement{
				Kind: quitStatement,
			})
		}
	}

	return statements, nil
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

	return cursor
}

func parseSingleKeywordStatement(tokens []*token, cursor uint, keyword keyword) (uint, bool) {
	if !tokens[cursor].equals(tokenFromKeyword(keyword)) {
		return cursor, false
	}

	cursor++
	cursor = eatAllToNextNewLine(tokens, cursor)

	return cursor, false
}

func eatNonKeywords(tokens []*token, cursor uint) uint {
	for uint(len(tokens)) < cursor && tokens[cursor].kind != keywordKind {
		cursor++
	}

	return cursor
}
