package uci

import (
	"testing"
)

func TestSetOption(t *testing.T) {

	expectations := []struct {
		source string
		name   string
		value  string
	}{
		{
			source: "setoption name foo value bar\n",
			name:   "foo",
			value:  "bar",
		},
		{
			source: "setoption name foo\n",
			name:   "foo",
			value:  "",
		},
		{
			source: "setoption NaME NalimovPath VALUE c:\\chess\\tb\\4;c:\\chess\\tb\\5\n",
			name:   "NalimovPath",
			value:  "c:\\chess\\tb\\4;c:\\chess\\tb\\5",
		},
	}

	for _, expectation := range expectations {
		statements, err := Parse(expectation.source)
		if err != nil {
			t.Fatal("error parsing test source", err)
		}

		if len(statements) != 1 {
			t.Errorf("expected 1 statement, but got %d", len(statements))
		}

		stmt := statements[0]

		if stmt.Kind != SetOptionStatementKind {
			t.Fatal("expected setOptionStatement, but got", stmt.Kind)
		}

		if stmt.SetOption.Name != expectation.name {
			t.Error("expected Name to be", expectation.name, ", but got", stmt.SetOption.Name)
		}

		if stmt.SetOption.Value != expectation.value {
			t.Error("expected Value to be", expectation.value, "but got", stmt.SetOption.Value)
		}
	}
}

func TestPosition(t *testing.T) {
	expectations := []struct {
		source     string
		kinds      []GoKind
		fenString  string
		isFen      bool
		isStartPos bool
		moves      []string
		wtime	   int
	}{
		{
			source:     "go searchmoves e1f1\n",
			kinds: 		[]GoKind{Go_searchMovesKind},
			fenString:  "0000",
			isFen:      true,
			isStartPos: false,
			moves:      []string{"e1f1"},
			wtime:		0,
		},
		{
			source:     "go wtime 1 searchmoves e1f1\n",
			kinds: 		[]GoKind{Go_searchMovesKind, Go_wtimeKind},
			fenString:  "0000",
			isFen:      true,
			isStartPos: false,
			moves:      []string{"e1f1"},
			wtime:		1,
		},
	}

	for _, expectation := range expectations {
		statements, err := Parse(expectation.source)
		if err != nil {
			t.Fatal("error parsing test source", err)
		}
		if len(statements) != 1 {
			t.Errorf("expected 1 statement, but got %d", len(statements))
		}
		stmnt := statements[0]
		if stmnt.Kind != GoStatementKind {
			t.Fatal("Expected go kind")
		}

		if len(stmnt.Go.Kinds) != len(expectation.kinds) {
			t.Fatal("Expected", len(expectation.kinds), "kinds, bot got", len(stmnt.Go.Kinds))
		}

		for _, actualKind := range stmnt.Go.Kinds {
			if !inKinds(actualKind, stmnt.Go.Kinds) {
				t.Fatal("Did not expected kind", actualKind)
			}
		}
	}
}

func inKinds(kind GoKind, list []GoKind) bool {
	for _, kindInList := range list {
		if kind == kindInList {
			return true
		}
	}

	return false
}
