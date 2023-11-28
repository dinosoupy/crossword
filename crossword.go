// crossword.go contains functions to load a crossword structure file, compute variables, overlaps and Neighbors
package crossword

import (
	"os"
	"strings"
)

const (
	ACROSS = 0
	DOWN   = 1
)

type Point struct {
	x, y int
}

type Variable struct {
	Coordinates Point           // row, col where the word starts
	Direction   int             // Direction of word where 0: across; 1: down
	Length      int             // Length of the word
	Cells       []Point         // Coordinates of all points
	Domain      []string        // store all possible values that the word can take
	Neighbors   map[Point]Point // store a map neighboring variables' coordinates and the point at which they overlap
}

var (
	Structure [][]bool
	Variables []*Variable
)

// NewCrossword parses a crossword structure file into a 2D slice of booleans
func NewCrossword(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	lines := strings.Split(string(data), "\n")
	height := len(lines)
	width := len(lines[0])
	Structure = make([][]bool, height)

	for row := 0; row < height; row++ {
		Structure[row] = make([]bool, width)
		Cells := strings.Split(lines[row], "")
		for col := 0; col < width; col++ {
			if Cells[col] == "_" {
				Structure[row][col] = true
			} else {
				Structure[row][col] = false
			}
		}
	}
	return nil
}

// ComputeVariables records the crossword Variables from a given structure slice
func ComputeVariables(structure [][]bool) {
	height, width := len(structure), len(structure[0])
	structureCopy := copy2DSlice(structure)

	for row := 0; row < height; row++ {
		for col := 0; col < width; col++ {
			// Across words
			Cells := make([]Point, 0)
			if structureCopy[row][col] {
				Cells = append(Cells, Point{row, col})
				Length := 1
				for i := col + 1; i < width; i++ {
					if structureCopy[row][i] {
						Length++
						Cells = append(Cells, Point{row, i})
						structureCopy[row][i] = false
					} else {
						break
					}
				}
				if Length > 1 {
					variable := &Variable{
						Coordinates: Point{row, col},
						Direction:   0,
						Length:      Length,
						Cells:       Cells,
					}
					variable.InitDomain()
					Variables = append(Variables, variable)
				}
			}

			// Down words
			Cells = make([]Point, 0)
			if structure[row][col] {
				Cells = append(Cells, Point{row, col})
				Length := 1
				for j := row + 1; j < height; j++ {
					if structure[j][col] {
						Length++
						Cells = append(Cells, Point{j, col})
						structure[j][col] = false
					} else {
						break
					}
				}
				if Length > 1 {
					variable := &Variable{
						Coordinates: Point{row, col},
						Direction:   1,
						Length:      Length,
						Cells:       Cells,
					}
					variable.InitDomain()
					Variables = append(Variables, variable)
				}

			}

		}
	}
}

// Overlaps returns the point where v1 overlaps with v2.
// Returns nil if v1 and v2 are not Neighbors
// TODO: Improve O(n^2) linear search
func (v1 *Variable) Overlaps(v2 *Variable) *Point {
	if v1.Coordinates != v2.Coordinates {
		for _, cellv1 := range v1.Cells {
			for _, cellv2 := range v2.Cells {
				if cellv1 == cellv2 {
					return &cellv2
				}
			}
		}
	}
	return nil
}

// InitNeighbors adds variables which overlap with v1 to `v1.Neighbors`
func (v1 *Variable) InitNeighbors() {
	v1.Neighbors = make(map[Point]Point)
	for _, v := range Variables {
		if overlap := v1.Overlaps(v); overlap != nil {
			v1.Neighbors[v.Coordinates] = *overlap
		}
	}
}

// InitDomain sets the Domain of a Variable to a slice of all words
// with Length = v.Length. This ensures variables are node consistent at
// the time of initialization
func (v *Variable) InitDomain() {
	v.Domain = make([]string, len(AllWords[v.Length]))
	copy(v.Domain, AllWords[v.Length])
}

func (v1 *Variable) Equals(v2 *Variable) bool {
	// Compare Coordinates
	if v1.Coordinates != v2.Coordinates {
		return false
	}

	// Compare Direction
	if v1.Direction != v2.Direction {
		return false
	}

	// Compare Length
	if v1.Length != v2.Length {
		return false
	}

	// Compare Cells
	if len(v1.Cells) != len(v2.Cells) {
		return false
	}
	for i := range v1.Cells {
		if v1.Cells[i] != v2.Cells[i] {
			return false
		}
	}

	// Compare Domain
	if len(v1.Domain) != len(v2.Domain) {
		return false
	}
	for i := range v1.Domain {
		if v1.Domain[i] != v2.Domain[i] {
			return false
		}
	}

	// Compare Neighbors
	if len(v1.Neighbors) != len(v2.Neighbors) {
		return false
	}
	for key, value := range v1.Neighbors {
		if otherValue, ok := v2.Neighbors[key]; !ok || value != otherValue {
			return false
		}
	}

	// All fields are equal
	return true
}

// GetVariable returns the Variable which starts ar a given point p
func GetVariable(p Point) Variable {
	for _, v := range Variables {
		if v.Coordinates == p {
			return *v
		}
	}
	return Variable{}
}

// Returns the index at which the overlap point occurs in the variable word
func (v *Variable) IndexOfCell(overlap Point) int {
	for i, cell := range v.Cells {
		if cell == overlap {
			return i
		}
	}
	return -1
}

// Helper function to deep copy structure slice
func copy2DSlice(s [][]bool) [][]bool {
	var copy = make([][]bool, len(s))
	for i, row := range s {
		copy[i] = make([]bool, len(row))
		for j, col := range row {
			copy[i][j] = col
		}
	}
	return copy
}
