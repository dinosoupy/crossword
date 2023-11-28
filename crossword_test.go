package crossword

import (
	"testing"
)

func TestNewCrossword(t *testing.T) {
	savedStructure, savedVariables := Structure, Variables
	defer func() {
		Structure, Variables = savedStructure, savedVariables
	}()

	NewCrossword("./structures/structure.txt")
	want := [][]bool{
		{false, true, true, true, false},
		{true, false, true, true, true},
		{true, false, true, true, false},
		{true, true, true, false, false},
		{false, false, true, false, false},
	}

	for i := 0; i < len(Structure); i++ {
		for j := 0; j < len(Structure[i]); j++ {
			if Structure[i][j] != want[i][j] {
				t.Errorf("NewCrossword(\"./structures/structure.txt\")\nGot:\t%v\nWant:\t%v", Structure, want)
			}
		}
	}
}

func TestVariables(t *testing.T) {
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

	want := []*Variable{
		{
			Coordinates: Point{0, 1},
			Direction:   0,
			Length:      3,
			Cells:       []Point{{0, 1}, {0, 2}, {0, 3}},
			Domain:      []string{"cat", "bat", "rat", "tab"},
			Neighbors:   map[Point]Point{{0, 2}: {0, 2}, {0, 3}: {0, 3}},
		},
		{
			Coordinates: Point{0, 2},
			Direction:   1,
			Length:      5,
			Cells:       []Point{{0, 2}, {1, 2}, {2, 2}, {3, 2}, {4, 2}},
			Domain:      []string{"arabs"},
			Neighbors:   map[Point]Point{{0, 1}: {0, 2}, {1, 2}: {1, 2}, {2, 2}: {2, 2}, {3, 0}: {3, 2}},
		},
		{
			Coordinates: Point{0, 3},
			Direction:   1,
			Length:      3,
			Cells:       []Point{{0, 3}, {1, 3}, {2, 3}},
			Domain:      []string{"cat", "bat", "rat", "tab"},
			Neighbors:   map[Point]Point{{0, 1}: {0, 3}, {1, 2}: {1, 3}, {2, 2}: {2, 3}},
		},
		{
			Coordinates: Point{1, 0},
			Direction:   1,
			Length:      3,
			Cells:       []Point{{1, 0}, {2, 0}, {3, 0}},
			Domain:      []string{"cat", "bat", "rat", "tab"},
			Neighbors:   map[Point]Point{{3, 0}: {3, 0}},
		},
		{
			Coordinates: Point{1, 2},
			Direction:   0,
			Length:      3,
			Cells:       []Point{{1, 2}, {1, 3}, {1, 4}},
			Domain:      []string{"cat", "bat", "rat", "tab"},
			Neighbors:   map[Point]Point{{0, 2}: {1, 2}, {0, 3}: {1, 3}},
		},
		{
			Coordinates: Point{2, 2},
			Direction:   0,
			Length:      2,
			Cells:       []Point{{2, 2}, {2, 3}},
			Domain:      []string{"at"},
			Neighbors:   map[Point]Point{{0, 2}: {2, 2}, {0, 3}: {2, 3}},
		},
		{
			Coordinates: Point{3, 0},
			Direction:   0,
			Length:      3,
			Cells:       []Point{{3, 0}, {3, 1}, {3, 2}},
			Domain:      []string{"cat", "bat", "rat", "tab"},
			Neighbors:   map[Point]Point{{0, 2}: {3, 2}, {1, 0}: {3, 0}},
		},
	}

	for i, v := range Variables {
		if !want[i].Equals(v) {
			t.Errorf("Variable mismatch:\nGot:\t%v\nWant:\t%v", v, want[i])
		}
	}
}
