package board

import "testing"

func TestPositionFromFileAndRank(t *testing.T) {
	tests := []struct {
		desc        string
		file        File
		rank        int
		expectation int
	}{
		{
			desc:        "A1",
			file:        A,
			rank:        1,
			expectation: 0,
		},
		{
			desc:        "B1",
			file:        B,
			rank:        1,
			expectation: 1,
		},
		{
			desc:        "A2",
			file:        A,
			rank:        2,
			expectation: 8,
		},
		{
			desc:        "G1",
			file:        G,
			rank:        1,
			expectation: 6,
		},
	}

	for _, tst := range tests {
		actual := indexFromFileAndRank(tst.file, tst.rank)
		if actual != tst.expectation {
			t.Errorf("expected %d from %s, but got %d", tst.expectation, tst.desc, actual)
		}
	}
}

func TestPositionFromIndex(t *testing.T) {
	for _, testCase := range []struct{pos int; expected Position}{
		{
			pos:      0,
			expected: Position{A, 1},
		},
		{
			pos:      8,
			expected: Position{A, 2},
		},
		{
			pos:      15,
			expected: Position{H, 2},
		},

	} {
		actual := PositionFromIndex(testCase.pos)
		if !actual.SameAs(&testCase.expected) {
			t.Errorf("expected %+v but got %+v", testCase.expected, actual)
		}
	}
}

func TestPosition120FromFileAndRank(t *testing.T) {
	tests := []struct {
		desc        string
		file        File
		rank        int
		expectation int
	}{
		{
			desc:        "A1",
			file:        A,
			rank:        1,
			expectation: 21,
		},
		{
			desc:        "B1",
			file:        B,
			rank:        1,
			expectation: 22,
		},
		{
			desc:        "A2",
			file:        A,
			rank:        2,
			expectation: 31,
		},
		{
			desc:        "G1",
			file:        G,
			rank:        1,
			expectation: 27,
		},
	}

	for _, tst := range tests {
		actual := index120FromFileAndRank(tst.file, tst.rank)
		if actual != tst.expectation {
			t.Errorf("expected %d from %s, but got %d", tst.expectation, tst.desc, actual)
		}
	}
}
