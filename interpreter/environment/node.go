package environment

type Node interface {
	// Evaluates the node in a certain execution environment
	Eval(*Environment) any
	// Gets the identifiers that a node references
	References() []string
}
