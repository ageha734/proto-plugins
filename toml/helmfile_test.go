package main

import (
	"testing"
)

func TestHelmfile(t *testing.T) {
	Run(TestConfig{
		Name: "helmfile",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("helmfile --version")
			return nil
		},
	})(t)
}
