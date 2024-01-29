package nodes

import "main/interpreter/environment"

type KeyAccess struct {
  Object Node
  Key string
}

func (ka *KeyAccess) Eval(env *environment.Environment) any {
  return ka.Object.Eval(env).(map[string]any)[ka.Key]
} 