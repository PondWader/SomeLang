package interpreter

import "main/interpreter/environment"

type TypeEnvironment struct {
	environment.Environment
	ReturnType TypeDef
	Returned   bool
	parent     *TypeEnvironment
}

func NewTypeEnvironment(parent *TypeEnvironment, returnType TypeDef) *TypeEnvironment {
	var parentEnv *environment.Environment
	if parent != nil {
		parentEnv = &parent.Environment
	}
	return &TypeEnvironment{*environment.New(parentEnv, nil, environment.Call{}), returnType, false, parent}
}

func (e *TypeEnvironment) NewChild(returnType TypeDef) *TypeEnvironment {
	return NewTypeEnvironment(e, returnType)
}

func (e *TypeEnvironment) GetReturnType() TypeDef {
	if e.ReturnType != nil {
		return e.ReturnType
	} else if e.parent != nil {
		return e.parent.GetReturnType()
	}
	return nil
}

func (e *TypeEnvironment) GetParent() *TypeEnvironment {
	return e.parent
}
