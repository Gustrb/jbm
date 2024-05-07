package cli

import (
	"errors"
	"fmt"

	"github.com/Gustrb/jbm/src/core"
)

type CLI struct {
	ProgramName string
	Arguments   []string
}

func CreateCLI(args []string) (*CLI, error) {
	if len(args) < 1 {
		return nil, errors.New("No arguments provided")
	}

	return &CLI{
		ProgramName: args[0],
		Arguments:   args[1:],
	}, nil
}

func (cli *CLI) DumpUsage() {
	progname := cli.ProgramName
	fmt.Printf("Usage: %s [options] <mainclass> [args...]\n\t\t", progname)
	fmt.Println("(to execute a class)")

	fmt.Printf("\tor %s [options] -jar <jarfile> [args...]\n\t\t", progname)
	fmt.Println("(to execute a jar file)")

	fmt.Printf("\tor %s [options] -m <module>[/<mainclass>] [args...]\n\t", progname)
	fmt.Printf("\t%s [options] --module <module>[/<mainclass>] [args...]\n\t\t", progname)
	fmt.Printf("(to execute the main class in a module)\n")

	fmt.Println(" The arguments after the main class, -jar <jarfile>, -m or --module")
	fmt.Println(" <module>/<mainclass> are specified as the arguments for the main class.")
}

func (cli *CLI) validateArguments() error {
	return nil
}

func (cli *CLI) Run() error {
	if len(cli.Arguments) < 1 {
		cli.DumpUsage()
		return errors.New("No arguments provided")
	}

	if err := cli.validateArguments(); err != nil {
		return err
	}

	if err := core.RunJBM(cli.Arguments); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
