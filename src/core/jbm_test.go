package core_test

import (
	"testing"

	"github.com/Gustrb/jbm/src/core"
)

func TestShouldFailIfNoFilepathIsProvided(t *testing.T) {
	args := []string{"--jar"}
	err := core.RunJBM(args)

	if err == nil {
		t.Error("expected an error, got nil")
	}
}

func TestShouldReturnJarTypeIfJarFlagIsPassed(t *testing.T) {
	args := []string{"--jar", "file.jar"}
	err := core.RunJBM(args)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestShouldReturnModuleTypeIfModuleFlagIsPassed(t *testing.T) {
	args := []string{"-m", "module-info.class"}
	err := core.RunJBM(args)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}
