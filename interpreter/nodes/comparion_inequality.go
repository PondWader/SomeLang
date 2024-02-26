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

type InequalityComparison[T int8 | int16 | int32 | int64 | uint8 | uint16 | uint32 | uint64 | float32 | float64] struct {
	Type      ComparisonType
	LeftSide  environment.Node
	RightSide environment.Node
}

func (c *InequalityComparison[T]) Eval(env *environment.Environment) any {
	if c.Type == ComparisonEquals {
		return c.LeftSide.Eval(env) == c.RightSide.Eval(env)
	} else if c.Type == ComparisonGreaterThan {
		return c.LeftSide.Eval(env).(T) > c.RightSide.Eval(env).(T)
	} else if c.Type == ComparisonGreaterThanOrEquals {
		return c.LeftSide.Eval(env).(T) >= c.RightSide.Eval(env).(T)
	} else if c.Type == ComparisonLessThan {
		return c.LeftSide.Eval(env).(T) < c.RightSide.Eval(env).(T)
	} else if c.Type == ComparisonLessThanOrEquals {
		return c.LeftSide.Eval(env).(T) <= c.RightSide.Eval(env).(T)
	}
	return nil
}

func (c *InequalityComparison[T]) References() []string {
	return append(c.RightSide.References(), c.LeftSide.References()...)
}
