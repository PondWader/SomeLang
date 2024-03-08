package nodes

import "main/interpreter/environment"

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
