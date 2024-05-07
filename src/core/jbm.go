package core

import (
	"bytes"
	"errors"

	"github.com/Gustrb/jbm/src/utils"
)

var Flags = []string{"--jar", "-m"}

type ExecutionContext struct {
	Filepath string
	Type     string
}

// Run executes the file that was passed as an argument.
//
// It first reads the file, parses the bytecode, and then executes it.
func (ctx *ExecutionContext) Run() error {
	fcontent, err := utils.ReadFileContent(ctx.Filepath)
	if err != nil {
		return err
	}

	// create a reader from the file content
	reader := bytes.NewReader(fcontent)
	switch ctx.Type {
	case "class":
		err = ExecuteClassFile(reader)
	}

	return err
}

// findExecutionType returns the type of the file that is being executed.
//
// If the --jar flag is passed, the type is "jar".
// If the -m flag is passed, the type is "module".
// If there is no type flag passed, the type is "class".
func findExecutionType(args []string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == "--jar" {
			return "jar"
		} else if args[i] == "-m" {
			return "module"
		}
	}

	return "class"
}

// findFilepath returns the filepath of the file that was used to execute the program.
//
// We can assume that the filepath is the argument that is not a flag.
func findFilepath(args []string) string {
	for i := 0; i < len(args); i++ {
		// I guess we couldve used a map here, but since the list is so small, it's not worth it.
		// maybe sometime in the future.
		if !utils.Contains(Flags, args[i]) {
			return args[i]
		}
	}

	return ""
}

func RunJBM(args []string) error {
	filepath := findFilepath(args)

	if filepath == "" {
		return errors.New("no filepath provided")
	}

	executionType := findExecutionType(args)

	ctx := ExecutionContext{
		Filepath: filepath,
		Type:     executionType,
	}

	err := ctx.Run()

	return err
}
