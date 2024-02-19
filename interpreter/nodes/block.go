package nodes

import "main/interpreter/environment"

type Block struct {
	Nodes []environment.Node
}

func (b *Block) Eval(env *environment.Environment) any {
	for _, node := range b.Nodes {
		node.Eval(env)
	}
	return nil
}

func (b *Block) References() []string {
	refs := make([]string, 0)
	for _, node := range b.Nodes {
		refs = append(refs, node.References()...)
	}
	return refs
}
