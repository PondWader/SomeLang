package environment

import (
	"fmt"
	"os"
	"strconv"
)

// Execution environment handles the storing of values, garbage collection, and evaluation of the AST
type Environment struct {
	identifiers map[string]any
	parent      *Environment
	// call stores which function call initialized the environment
	Call Call

	returnCallback func(any)
	// Public value of whether or not loops are broken and should exit (e.g. after a return or break statement)
	IsBroken bool

	ast          []Node
	position     int
	attachedRefs map[string][]string

	modules map[string]map[string]any
}

type Call struct {
	File         string
	Line         int
	FunctionName string
}

func New(parent *Environment, call Call, modules map[string]map[string]any) *Environment {
	return &Environment{
		identifiers:  make(map[string]any),
		parent:       parent,
		Call:         call,
		attachedRefs: make(map[string][]string),
		modules:      modules,
	}
}

func (e *Environment) Get(name string) any {
	if value, ok := e.identifiers[name]; ok {
		return value
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	return nil
}

func (e *Environment) Set(name string, value any) {
	e.identifiers[name] = value
}

// Sets a value in a parent environment with a depth of how many parent environments to go back
func (env *Environment) SetWithDepth(name string, value any, depth int) {
	for i := 0; i < depth; i++ {
		env = env.GetParent()
		if env == nil {
			panic("Depth is greater than total available depth")
		}
	}
	env.Set(name, value)
}

func (e *Environment) NewChild(call Call) *Environment {
	child := New(e, call, e.modules)
	child.SetReturnCallback(e.returnCallback)
	return child
}

func (e *Environment) GetParent() *Environment {
	return e.parent
}

func (e *Environment) GetCallStackOutput() string {
	output := "\tFile: " + e.Call.File + ", Line: " + strconv.Itoa(e.Call.Line) + ", In " + e.Call.FunctionName
	if e.parent != nil {
		output += "\n" + e.parent.GetCallStackOutput()
	}
	return output
}

func (e *Environment) Execute(ast []Node) {
	e.ast = ast
	for i, node := range ast {
		if e.IsBroken {
			return
		}
		e.position = i
		node.Eval(e)
		e.RunGC()
	}
}

func (e *Environment) SetReturnCallback(cb func(v any)) {
	e.returnCallback = func(v any) {
		e.IsBroken = true
		cb(v)
	}
}

func (e *Environment) Return(v any) {
	e.returnCallback(v)
}

// Declares references attached to a certain identifier.
// This allows the garbage collector to mark identifiers relied on by another identifier as in use.
func (e *Environment) AttachReferences(name string, refs []string) {
	e.attachedRefs[name] = refs
}

func (e *Environment) Panic(msg ...any) {
	fmt.Println(append([]any{"panic:"}, msg...)...)
	fmt.Println(e.GetCallStackOutput())
	os.Exit(1)
}

func (e *Environment) GetBuiltInModule(module string) map[string]any {
	return e.modules[module]
}
