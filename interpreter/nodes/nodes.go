package nodes

import "main/interpreter/environment"

type Node interface {
  Eval(*environment.Environment) any
}