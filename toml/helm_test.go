package main

import (
	"testing"
)

func TestHelm(t *testing.T) {
	Run(TestConfig{
		Name: "helm",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("helm version")
			return nil
		},
	})(t)
}
