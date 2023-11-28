package crossword

import (
	"fmt"
)

func main() {
	NewCrossword("./structures/structure.txt")
	ComputeVariables(Structure)
	emptyAssignment := make(map[*Variable]string)
	fmt.Println(BacktrackAC3(emptyAssignment))
}
