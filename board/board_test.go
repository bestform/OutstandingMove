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
		actual := positionFromFileAndRank(tst.file, tst.rank)
		if actual != tst.expectation {
			t.Errorf("expected %d from %s, but got %d", tst.expectation, tst.desc, actual)
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
		actual := position120FromFileAndRank(tst.file, tst.rank)
		if actual != tst.expectation {
			t.Errorf("expected %d from %s, but got %d", tst.expectation, tst.desc, actual)
		}
	}
}
