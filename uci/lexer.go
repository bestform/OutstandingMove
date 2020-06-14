package uci

import (
	"fmt"
	"sort"
	"strings"
)

type keyword string

const (
	uci         keyword = "uci"
	debug       keyword = "debug"
	on          keyword = "on"
	off         keyword = "off"
	isready     keyword = "isready"
	setoption   keyword = "setoption"
	register    keyword = "register"
	later       keyword = "later"
	name        keyword = "name"
	code        keyword = "code"
	value       keyword = "value"
	ucinewgame  keyword = "ucinewgame"
	position    keyword = "position"
	fen         keyword = "fen"
	startpos    keyword = "startpos"
	moves       keyword = "moves"
	go_         keyword = "go"
	searchmoves keyword = "searchmoves"
	ponder      keyword = "ponder"
	wtime       keyword = "wtime"
	btime       keyword = "btime"
	winc        keyword = "winc"
	binc        keyword = "binc"
	movestogo   keyword = "movestogo"
	depth       keyword = "depth"
	nodes       keyword = "nodes"
	mate        keyword = "mate"
	movetime    keyword = "movetime"
	infinite    keyword = "infinite"
	stop        keyword = "stop"
	ponderhit   keyword = "ponderhit"
	quit        keyword = "quit"
)

type symbol byte

const (
	newLine symbol = '\n'
)

type tokenKind uint

const (
	stringKind tokenKind = iota
	symbolKind
	keywordKind
	longAlgebraicNotation
)

type token struct {
	kind  tokenKind
	value string
}

func (t *token) equals(other *token) bool {
	return t.kind == other.kind && t.value == other.value
}

type cursor uint

type lexer func(source string, cur cursor) (*token, cursor, bool)

func Lex(source string) ([]*token, error) {
	var cur cursor

	lexers := []lexer{symbolLexer, keywordLexer, longAlgebraicNotationLexer, stringLexer}
	var tokens []*token
lex:
	for cur < cursor(len(source)) {
		cur = eatWhitespace(source, cur)
		for _, lexer := range lexers {
			t, newCursor, ok := lexer(source, cur)
			if ok {
				tokens = append(tokens, t)
				cur = newCursor
				continue lex
			}
		}

		return nil, fmt.Errorf("error lexing at position %d", cur)
	}

	return tokens, nil
}

func expect(source string, ic cursor, searchString string) bool {
	return strings.HasPrefix(source[ic:], searchString)
}

func stringLexer(source string, ic cursor) (*token, cursor, bool){
	cur := ic
	value := ""
	for cur < cursor(len(source)) && source[cur] != ' ' && source[cur] != uint8(newLine) {
		value += string(source[cur])
		cur++
	}

	return &token{
		kind:  stringKind,
		value: value,
	}, cur, true
}

func longAlgebraicNotationLexer(source string, ic cursor) (*token, cursor, bool) {
	if len(source[ic:]) < 4 {
		return nil, ic, false
	}

	if expect(source, ic, "0000") {
		return &token{
			kind:  longAlgebraicNotation,
			value: "0000",
		}, ic+4, true
	}

	validFiles := []byte{'a','b','c','d','e','f','g','h'}
	validPiecesForPromotion := []byte{'r','k','q','n','b'}

	cur := ic
	value := ""

	for range []int{0, 1} {
		current := source[cur]
		isValidFile := false
		for _, f := range validFiles {
			if current == f {
				isValidFile = true
				break
			}
		}
		if !isValidFile {
			return nil, ic, false
		}
		value += string(current)
		cur++

		current = source[cur]

		if current > '8' || current < '1' {
			return nil, ic, false
		}
		value += string(current)
		cur++
	}

	// promotion?
	if cursor(len(source)) > cur {
		current := source[cur]
		if current != ' ' && current != '\n' {
			isValidPiece := false
			for _, p := range validPiecesForPromotion {
				if current == p {
					isValidPiece = true
					break
				}
			}
			if !isValidPiece {
				return nil, ic, false
			}
			value += string(current)
			cur++
		}
	}

	return &token{
		kind:  longAlgebraicNotation,
		value: value,
	}, cur, true
}

func symbolLexer(source string, ic cursor) (*token, cursor, bool) {
	cur := ic

	current := source[ic]

	for _, sym := range []symbol{newLine} {
		if current == byte(sym) {
			return &token{
				kind:  symbolKind,
				value: string(sym),
			}, cur + 1, true
		}
	}

	return nil, ic, false
}

func keywordLexer(incSource string, ic cursor) (*token, cursor, bool) {
	cur := ic

	source := strings.ToLower(incSource)

	keywords := []keyword{value, binc, btime, code, debug, depth, fen, go_, infinite, isready, later, mate, moves, movestogo, movetime, name, nodes, off, on, ponder, ponderhit, position, quit, register, searchmoves, setoption, startpos, stop, uci, ucinewgame, winc, wtime}

	sort.Slice(keywords, func(a int, b int) bool {
		return len(string(keywords[a])) > len(string(keywords[b]))
	})

	for _, kw := range keywords {
		if strings.HasPrefix(source[cur:], string(kw)) {
			return &token{
				kind:  keywordKind,
				value: string(kw),
			}, cur + cursor(len(string(kw))), true
		}
	}

	return nil, ic, false
}

func eatWhitespace(source string, ic cursor) cursor {
	cur := ic

	for cur < cursor(len(source)) && source[cur] == ' ' {
		cur++
	}

	return cur
}
