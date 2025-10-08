package main

import (
	"testing"
)

func TestTerragrunt(t *testing.T) {
	Run(TestConfig{
		Name: "terragrunt",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("terragrunt --version")
			return nil
		},
	})(t)
}
