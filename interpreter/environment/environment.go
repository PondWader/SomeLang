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
	// return can store a callback that returns a value in a function call
	Return func(any)

	ast          []Node
	position     int
	attachedRefs map[string][]string
}

type Call struct {
	File         string
	Line         int
	FunctionName string
}

func New(parent *Environment, ast []Node, call Call) *Environment {
	return &Environment{
		identifiers:  make(map[string]any),
		parent:       parent,
		Call:         call,
		ast:          ast,
		attachedRefs: make(map[string][]string),
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
	child := New(e, nil, call)
	child.Return = e.Return
	return child
}

func (e *Environment) GetParent() *Environment {
	return e.parent
}

func (e *Environment) GetCallStackOutput() string {
	output := "File, " + e.Call.File + ", Line, " + strconv.Itoa(e.Call.Line) + ", In " + e.Call.FunctionName
	if e.parent != nil {
		output += "\n" + e.parent.GetCallStackOutput()
	}
	return output
}

func (e *Environment) Execute() {
	for i, node := range e.ast {
		e.position = i
		node.Eval(e)
		e.RunGC()
	}
}

// Declares references attached to a certain identifier.
// This allows the garbage collector to mark identifiers relied on by another identifier as in use.
func (e *Environment) AttachReferences(name string, refs []string) {
	e.attachedRefs[name] = refs
}

func (e *Environment) Panic(msg ...any) {
	fmt.Println("panic:", msg)
	// TODO: Stack trace
	os.Exit(1)
}
