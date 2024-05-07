package cli_test

import (
	"testing"

	"github.com/Gustrb/jbm/src/cli"
)

func TestItShouldRequireAtLeast2Arguments(t *testing.T) {
	empty := []string{}

	_, err := cli.CreateCLI(empty)
	if err == nil {
		t.Error("Expected error, got nil")
	}
}

func TestItShouldFindTheProgramName(t *testing.T) {
	arguments := []string{"jbm", "arg1", "arg2"}

	c, _ := cli.CreateCLI(arguments)

	if c.ProgramName != arguments[0] {
		t.Errorf("Expected 'jbm', got %s", c.ProgramName)
	}
}
