package environment

func (e *Environment) RunGC() {
	// Mark variables that are in use in a hash map
	inUse := make(map[string]struct{}, len(e.identifiers))
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
			delete(e.identifiers, ident)
			delete(e.attachedRefs, ident)
		}
	}
}
