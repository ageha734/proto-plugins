package main

import (
	"testing"
)

func TestShfmt(t *testing.T) {
	Run(TestConfig{
		Name: "shfmt",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("shfmt --version")
			return nil
		},
	})(t)
}
