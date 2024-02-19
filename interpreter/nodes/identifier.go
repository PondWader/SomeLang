package nodes

import (
	"main/interpreter/environment"
)

type Identifier struct {
	Name string
}

func (i *Identifier) Eval(env *environment.Environment) any {
	return env.Get(i.Name)
}

func (i *Identifier) References() []string {
	return []string{i.Name}
}
