package main

import (
	"testing"
)

func TestDprint(t *testing.T) {
	Run(TestConfig{
		Name: "dprint",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("dprint --version")
			return nil
		},
	})(t)
}
