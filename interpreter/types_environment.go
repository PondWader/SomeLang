package interpreter

type TypeEnvironment struct {
	identifiers map[string]TypeDef
	returnType  TypeDef
	returned    bool
	parent      *TypeEnvironment
	Depth       int
}

func NewTypeEnvironment(parent *TypeEnvironment, returnType TypeDef, depth int) *TypeEnvironment {
	return &TypeEnvironment{make(map[string]TypeDef), returnType, false, parent, depth}
}

// Creates a new type environment with the current instance as it's parent
func (e *TypeEnvironment) NewChild(returnType TypeDef) *TypeEnvironment {
	return NewTypeEnvironment(e, returnType, e.Depth+1)
}

func (e *TypeEnvironment) GetReturnType() TypeDef {
	if e.returnType != nil {
		return e.returnType
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

// Gets a value by it's name and also returns the depth of the parent environment it was retrieved from
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

func (e *TypeEnvironment) GetReturned() bool {
	return e.returned
}

// Marks the environment as having had a "return" statement ini t
func (e *TypeEnvironment) SetReturned() {
	e.returned = true
}
