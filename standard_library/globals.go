package standardlibrary

import (
	"bufio"
	"fmt"
	"main/interpreter"
	"os"
)

var LenDef = interpreter.FuncDef{
	GenericTypeDef: interpreter.GenericTypeDef{Type: interpreter.TypeFunc},
	Args: []interpreter.TypeDef{
		interpreter.GenericTypeDef{Type: interpreter.TypeString},
	},
	ReturnType: interpreter.GenericTypeDef{Type: interpreter.TypeInt64},
}

func Len[A any, V string | []A](v V) int64 {
	return int64(len(v))
}

var PrintDef = interpreter.FuncDef{
	GenericTypeDef: interpreter.GenericTypeDef{Type: interpreter.TypeFunc},
	Args: []interpreter.TypeDef{
		interpreter.GenericTypeDef{Type: interpreter.TypeAny},
	},
	Variadic:   true,
	ReturnType: nil,
}

func Print(args ...any) {
	fmt.Println(args...)
}

var InputDef = interpreter.FuncDef{
	GenericTypeDef: interpreter.GenericTypeDef{Type: interpreter.TypeFunc},
	Args:           make([]interpreter.TypeDef, 0),
	ReturnType:     interpreter.GenericTypeDef{Type: interpreter.TypeString},
}

func Input() string {
	reader := bufio.NewReader(os.Stdin)
	line, _, _ := reader.ReadLine()
	return string(line)
}
