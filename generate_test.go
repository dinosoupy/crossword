package crossword

import (
	"testing"
)

func TestMostConstrainedVariable(t *testing.T) {
	savedAllWords, savedStructure, savedVariables := AllWords, Structure, Variables
	defer func() {
		AllWords, Structure, Variables = savedAllWords, savedStructure, savedVariables
	}()

	AllWords = map[int][]string{
		2: {"at"},
		3: {"cat", "bat", "rat", "tab"},
		5: {"arabs"},
	}

	Structure = [][]bool{
		{false, true, true, true, false},
		{true, false, true, true, true},
		{true, false, true, true, false},
		{true, true, true, false, false},
		{false, false, true, false, false},
	}

	ComputeVariables(Structure)
	for _, v := range Variables {
		v.InitDomain()
		v.InitNeighbors()
	}

	want := Variable{
		Coordinates: Point{0, 2},
		Direction:   1,
		Length:      5,
		Cells:       []Point{{0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2}},
		Domain:      []string{"arabs"},
		Neighbors:   map[Point]Point{{0, 1}: {0, 2}, {1, 2}: {1, 2}, {2, 2}: {2, 2}, {3, 0}: {3, 2}},
	}

	got := MostConstrainedVariable()

	if !want.Equals(got) {
		t.Errorf("Variable mismatch:\nGot:\t%v\nWant:\t%v", got, want)
	}
}
