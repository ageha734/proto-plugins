package main

import (
	"testing"
)

func TestAct(t *testing.T) {
	Run(TestConfig{
		Name: "act",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("act --version")
			return nil
		},
	})(t)
}
