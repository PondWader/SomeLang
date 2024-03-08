package interpreter

type TypeEnvironment struct {
	// TODO: abstract environment in to environment and then execution environment which extends it and possibly come up with a better name for this
	identifiers map[string]TypeDef
	customTypes map[string]TypeDef
	ReturnType  TypeDef
	Returned    bool
	parent      *TypeEnvironment
	// How many parent environments there are
	Depth int
}

func NewTypeEnvironment(parent *TypeEnvironment, returnType TypeDef, depth int) *TypeEnvironment {
	return &TypeEnvironment{make(map[string]TypeDef), make(map[string]TypeDef), returnType, false, parent, depth}
}

func (e *TypeEnvironment) NewChild(returnType TypeDef) *TypeEnvironment {
	return NewTypeEnvironment(e, returnType, e.Depth+1)
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

func (e *TypeEnvironment) Get(name string) (TypeDef, int) {
	return e.getWithDepthCounter(name, 0)
}

func (e *TypeEnvironment) getWithDepthCounter(name string, depth int) (TypeDef, int) {
	if val, ok := e.identifiers[name]; ok {
		return val, depth
	}
	if e.parent != nil {
		return e.parent.getWithDepthCounter(name, depth+1)
	}
	return nil, -1
}

func (e *TypeEnvironment) Set(name string, value TypeDef) {
	e.identifiers[name] = value
}

func (e *TypeEnvironment) SetReturned() {
	e.Returned = true
}
