package main

import (
	"flag"
	"fmt"
	"main/interpreter"
	standardlibrary "main/standard_library"
	keyvalue "main/standard_library/key_value.go"
	"os"
	"path/filepath"
)

func main() {
	entryPoint := flag.String("run", "", "The entry point file to run")
	runProfiler := flag.Bool("profile", false, "If passed the program execution will be profiled")
	flag.Parse()

	fmt.Println(*entryPoint, *runProfiler)
	var err error
	if *entryPoint == "" {
		fmt.Println("You must specify an entrypoint with the -run flag")
		return
	}
	*entryPoint, err = filepath.Abs(*entryPoint)
	if err != nil {
		fmt.Println("Error resolving entry point path:", err)
		return
	}
	content, err := os.ReadFile(*entryPoint)
	if err != nil {
		fmt.Println("Error reading entry point file:", err)
		return
	}

	parser := interpreter.NewParser(string(content), *entryPoint, map[string]interpreter.TypeDef{
		"print": standardlibrary.PrintDef,
		"input": standardlibrary.InputDef,
	}, map[string]map[string]interpreter.TypeDef{
		"key_value": {
			"open": keyvalue.OpenDef,
		},
	})
	ast := parser.Parse()

	interpreter.Execute(ast, *entryPoint, map[string]any{
		"print": standardlibrary.Print,
		"input": standardlibrary.Input,
	}, map[string]map[string]any{
		"key_value": {
			"open": keyvalue.Open,
		},
	})
}
