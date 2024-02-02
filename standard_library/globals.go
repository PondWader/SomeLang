package standardlibrary

import (
	"bufio"
	"fmt"
	"main/interpreter"
	"os"
	"strings"
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

func Trim(v string) string {
	return strings.TrimSpace(v)
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

func Input() string {
	reader := bufio.NewReader(os.Stdin)
	line, _, _ := reader.ReadLine()
	return string(line)
}
