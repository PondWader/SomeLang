package nodes

import "main/interpreter/environment"

type Block struct {
	Nodes    []Node
}

func (b *Block) Eval(env *environment.Environment) any {
	for _, node := range b.Nodes {
    node.Eval(env)
  }
	return nil
}
