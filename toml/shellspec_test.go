package main

import (
	"testing"
)

func TestShellspec(t *testing.T) {
	Run(TestConfig{
		Name: "shellspec",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("shellspec --version")
			return nil
		},
	})(t)
}
