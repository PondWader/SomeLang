package main

import (
	"flag"
	"fmt"
	"main/interpreter"
	standardlibrary "main/standard_library"
	"os"
	"path/filepath"
)

func main() {
	entryPoint := flag.String("run", "", "The entry point file to run")
	runProfiler := flag.Bool("profile", false, "If passed the program execution will be profiled")
	flag.Parse()

	fmt.Println(*entryPoint, *runProfiler)
	var err error
	*entryPoint, err = filepath.Abs(*entryPoint)
	if err != nil {
		fmt.Println("Error resolving entry point path:", err)
	}
	content, err := os.ReadFile(*entryPoint)
	if err != nil {
		fmt.Println("Error reading entry point file:", err)
	}

	ast := interpreter.Parse(string(content), *entryPoint, map[string]interpreter.TypeDef{
		"print": standardlibrary.PrintDef,
		"input": standardlibrary.InputDef,
	})
	fmt.Println("Executing...")
	interpreter.Execute(ast, *entryPoint, map[string]any{
		"print": standardlibrary.Print,
		"input": standardlibrary.Input,
	})
}
