package main

import (
	"testing"
)

func TestShellcheck(t *testing.T) {
	Run(TestConfig{
		Name: "shellcheck",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("shellcheck --version")
			return nil
		},
	})(t)
}
