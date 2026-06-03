package main

import (
	"testing"
)

func TestSkaffold(t *testing.T) {
	Run(TestConfig{
		Name: "skaffold",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("skaffold version")
			return nil
		},
	})(t)
}
