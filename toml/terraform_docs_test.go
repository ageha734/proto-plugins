package main

import (
	"testing"
)

func TestTerraformDocs(t *testing.T) {
	Run(TestConfig{
		Name: "terraform-docs",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("terraform-docs --version")
			return nil
		},
	})(t)
}
