package nodes

import "main/interpreter/environment"

type VarDeclaration struct {
	Identifier string
	Value      Node
}

func (vd *VarDeclaration) Eval(env *environment.Environment) any {
  env.Set(vd.Identifier, vd.Value.Eval(env))
  return nil
}

func (vd *VarDeclaration) Type(env *environment.Environment) string {
  return ""
}