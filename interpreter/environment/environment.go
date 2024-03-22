package environment

import (
	"fmt"
	"main/profiler"
	"os"
	"strconv"
	"time"
)

// Execution environment handles the storing of values, garbage collection, and evaluation of the AST
type Environment struct {
	identifiers map[string]any
	parent      *Environment
	// A pointer to an address that is always updated to the current environment being executed
	currentExecutionEnv **Environment
	// call stores which function call initialized the environment
	Call Call

	returnCallback func(any)
	// Public value of whether or not loops are broken and should exit (e.g. after a return or break statement)
	IsBroken bool

	ast           []Node
	position      int
	attachedRefs  map[string][]string
	profile       bool
	profileResult *profiler.ProfileResult

	modules map[string]map[string]any
}

// A Call instance represents a function call in the call stack for stack trace displays
type Call struct {
	File string
	Line int
	Name string
}

func New(parent *Environment, call Call, modules map[string]map[string]any, profile bool) *Environment {
	var profileResult *profiler.ProfileResult
	if profile && call.Name != "" {
		profileResult = &profiler.ProfileResult{
			Name: call.Name,
			// Have to initialize the result array
			SubPrograms: make([]*profiler.ProfileResult, 0),
		}
	}

	var currentExecutionEnv **Environment
	if parent != nil {
		currentExecutionEnv = parent.currentExecutionEnv
	}

	return &Environment{
		identifiers:         make(map[string]any),
		parent:              parent,
		Call:                call,
		attachedRefs:        make(map[string][]string),
		modules:             modules,
		profile:             profile,
		profileResult:       profileResult,
		currentExecutionEnv: currentExecutionEnv,
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
	child := New(e, call, e.modules, e.profile)
	child.SetReturnCallback(e.returnCallback)
	return child
}

func (e *Environment) GetParent() *Environment {
	return e.parent
}

// Generates the call stack of the up to the current call
func (e *Environment) GetCallStackOutput() string {
	output := ""
	if e.Call.Name != "" {
		output = "\tFile: " + e.Call.File + ", Line: " + strconv.Itoa(e.Call.Line) + ", In " + e.Call.Name
	}
	if e.parent != nil {
		if output != "" {
			output += "\n"
		}
		output += e.parent.GetCallStackOutput()
	}
	return output
}

func (e *Environment) Execute(ast []Node) {
	var prevExecutionEnv *Environment
	// Update the current execution env
	if e.currentExecutionEnv != nil {
		prevExecutionEnv = *e.currentExecutionEnv
		*e.currentExecutionEnv = e
	} else {
		e.currentExecutionEnv = &e
	}
	// Store start time for profiling
	startTime := time.Now()

	e.ast = ast
	for i, node := range ast {
		if e.IsBroken {
			// If IsBroken is true, the environment should stop executing since a return statement has been reached
			return
		}
		e.position = i
		node.Eval(e)
		e.RunGC()
	}

	if e.profileResult != nil {
		// Measure the execution time
		e.profileResult.Duration = time.Since(startTime)
	}
	if prevExecutionEnv != nil {
		// Reset the previous executioon env
		*e.currentExecutionEnv = prevExecutionEnv
	}
}

// Sets a function that will be called when the environment returns
func (e *Environment) SetReturnCallback(cb func(v any)) {
	e.returnCallback = func(v any) {
		e.IsBroken = true
		cb(v)
	}
}

func (e *Environment) Return(v any) {
	e.returnCallback(v)
}

// Saves a profile result for a function call within the current environment
func (e *Environment) ProfileFunctionCall(result *profiler.ProfileResult) {
	if e.profileResult != nil {
		e.profileResult.SubPrograms = append(e.profileResult.SubPrograms, result)
	} else if e.profile && e.parent != nil {
		e.parent.ProfileFunctionCall(result)
	}
}

func (e *Environment) GetProfileResult() *profiler.ProfileResult {
	return e.profileResult
}

// Get's the environment anywhere in tree currently being executed
func (e *Environment) GetCurrentExecutionEnv() *Environment {
	return *e.currentExecutionEnv
}

// Declares references attached to a certain identifier.
// This allows the garbage collector to mark identifiers relied on by another identifier as in use.
func (e *Environment) AttachReferences(name string, refs []string) {
	e.attachedRefs[name] = refs
}

// Performs a runtime panic, throwing an error and exiting the program
func (e *Environment) Panic(msg ...any) {
	fmt.Println(append([]any{"panic:"}, msg...)...)
	fmt.Println(e.GetCallStackOutput())
	os.Exit(1)
}

func (e *Environment) GetBuiltInModule(module string) map[string]any {
	return e.modules[module]
}
