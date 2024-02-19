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
	LeftSide  environment.Node
	RightSide environment.Node
}

func (c *Comparison) Eval(env *environment.Environment) any {
	if c.Type == ComparisonEquals {
		return c.LeftSide.Eval(env) == c.RightSide.Eval(env)
	} else if c.Type == ComparisonGreaterThan {
		// c.LeftSide.Eval(env) > c.RightSide.Eval(env)
	}
	return nil
}

func (c *Comparison) References() []string {
	return append(c.RightSide.References(), c.LeftSide.References()...)
}
