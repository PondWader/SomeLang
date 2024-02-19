package environment

import "fmt"

type Node interface {
	// Evaluates the node in a certain execution environment
	Eval(*Environment) any
	// Gets the identifiers that a node references
	References() []string
}

func (e *Environment) RunGC() {
	fmt.Println("Running garbage collector")

	// Mark variables that are in use in a hash map
	inUse := make(map[string]struct{}, 0)
	for i := e.position + 1; i < len(e.ast); i++ {
		for _, ref := range e.ast[i].References() {
			inUse[ref] = struct{}{}
			if attachedRefs, ok := e.attachedRefs[ref]; ok {
				for _, ref := range attachedRefs {
					inUse[ref] = struct{}{}
				}
			}
		}
	}

	// Sweep identifiers that are out of scope from memory
	for ident := range e.identifiers {
		if _, ok := inUse[ident]; !ok {
			fmt.Println("Sweeped", ident)
			delete(e.identifiers, ident)
			delete(e.attachedRefs, ident)
		}
	}

	fmt.Println("Finished run of garbage collector")
}
