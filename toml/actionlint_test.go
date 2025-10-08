package main

import (
	"testing"
)

func TestActionlint(t *testing.T) {
	Run(TestConfig{
		Name: "actionlint",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("actionlint --version")
			return nil
		},
	})(t)
}
