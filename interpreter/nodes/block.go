package nodes

import "main/interpreter/environment"

// Node that represents a nested code block within the program such as a function body
type Block struct {
	Nodes []environment.Node
}

func (n *Block) Eval(env *environment.Environment) any {
	env.Execute(n.Nodes)
	return nil
}

func (n *Block) References() []string {
	refs := make([]string, 0)
	for _, node := range n.Nodes {
		refs = append(refs, node.References()...)
	}
	return refs
}
