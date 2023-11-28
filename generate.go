package crossword

// Returns a Variable that is not a part of Assignment.
func SelectUnassignedVariable(assignment map[*Variable]string) Variable {
	for _, v := range Variables {
		if _, ok := assignment[v]; !ok {
			return *v
		}
	}
	return Variable{}
}

// Returns true if all variables have been assigned
func AllAssignmentsComplete(assignment map[*Variable]string) bool {
	if len(Variables) == len(assignment) {
		for _, v := range Variables {
			if _, ok := assignment[v]; !ok {
				return false
			}
		}
		return true
	} else {
		return false
	}
}

// EnforceArcConsistency(var1, var2) reduces var1's domain such that it is
// arc consistent wrt var2. var1 is arc consistent with var2 if for each value in its
// domain, there exists a value of var2 for which the overlap between var and var2
// has the same letter
func EnforceArcConsistency(var1, var2 *Variable, overlap Point) {
	reducedDomain := make([]string, 0)
	for _, val1 := range var1.Domain {
		// find the char at the overlap point
		charAtOverlap := []rune(val1)[var1.IndexOfCell(overlap)]
		possibleVal2s := Lexicon[var2.Length][var2.IndexOfCell(overlap)].bitarrays[charAtOverlap]

		if !possibleVal2s.IsEmpty() {
			reducedDomain = append(reducedDomain, val1)
		}
	}
	var1.Domain = reducedDomain
}

// AC3 starts with a variable and appends its neighbors to the Queue
// Until Queue is empty, the first var is made to be arc consistent with its neighbors.
// Each neighbor is enqueued to Queue. If at any point, var1's domain is exhausted, return false.
// If all variables are arc consistent by the end, return true
func AC3(queue []Variable) bool {
	for len(queue) > 0 {
		var1 := Dequeue(queue)
		// Iterate through all neighbors and enforce arc consistency
		for var2Coords, overlap := range var1.Neighbors {
			var2 := GetVariable(var2Coords)
			Enqueue(queue, var2)
			EnforceArcConsistency(&var1, &var2, overlap)
			if len(var1.Domain) == 0 {
				return false
			}
		}
	}
	return true
}

func BacktrackAC3(assignment map[*Variable]string) map[*Variable]string {
	// base case
	if AllAssignmentsComplete(assignment) {
		return assignment
	}

	// select an unassigned variable from Variables
	var1 := SelectUnassignedVariable(assignment)
	savedDomain := var1.Domain
	// For now, values in domain are not ordered
	// Instead of LCV, loop through all values in domain
	for _, val := range var1.Domain {
		// set partial assignment
		assignment[&var1] = val
		// set domain to only contain val
		var1.Domain = []string{val}
		// make a queue called arcs which stores var1 initially
		arcs := []Variable{var1}
		// call AC3() on the first var, which loops through all neighbors
		// and calls AC3() on them
		AC3(arcs)
		result := BacktrackAC3(assignment)

		if result != nil {
			return result
		}

		// reset assignment and var1 domain
		delete(assignment, &var1)
		var1.Domain = savedDomain
	}
	return nil
}

func Enqueue(q []Variable, v Variable) []Variable {
	q = append(q, v)
	return q
}

func Dequeue(q []Variable) Variable {
	if len(q) > 0 { // check if queue is non empty
		v := q[0]
		q = q[1:]
		return v
	} else {
		return Variable{} // return value for underflow case
	}
}
