package nodes

import "main/interpreter/environment"

type ComparisonType uint8

const (
	ComparisonEquals ComparisonType = iota
	ComparisonGreaterThan
	ComparisonLessThan
	ComparisonGreaterThanOrEquals
	ComparisonLessThanOrEquals
)

type Comparison struct {
	Type      ComparisonType
	LeftSide  Node
	RightSide Node
}

func (c Comparison) Eval(env *environment.Environment) any {
	if c.Type == ComparisonEquals {
		return c.LeftSide.Eval(env) == c.RightSide.Eval(env)
	} else if c.Type == ComparisonGreaterThan {
		// c.LeftSide.Eval(env) > c.RightSide.Eval(env)
	}
	return nil
}
