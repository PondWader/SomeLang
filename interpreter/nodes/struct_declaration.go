package nodes

import (
	"main/interpreter/environment"
)

type StructDeclaration struct {
	Name    string
	Methods []environment.Node
}

// At runtime a struct is a function that can be called with the struct parameters to create a new instance
// A struct instance is an array of type any to store different data types
// All validation should have been done ahead of time by the parser

func (n *StructDeclaration) Eval(env *environment.Environment) any {
	env.Set(n.Name, func(properties ...any) {
		methodEnv := env.NewChild(env.Call)

		for _, method := range n.Methods {
			properties = append(properties, method.Eval(methodEnv))
		}

		methodEnv.Set("self", properties)
	})
	return nil
}

func (n *StructDeclaration) References() []string {
	refs := make([]string, 0)
	for _, method := range n.Methods {
		refs = append(refs, method.References()...)
	}
	return refs
}
