package main

import (
	"testing"
)

func TestKubectx(t *testing.T) {
	Run(TestConfig{
		Name: "kubectx",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("kubectx --version")
			return nil
		},
	})(t)
}
