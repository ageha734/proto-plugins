package main

import (
	"testing"
)

func TestKustomize(t *testing.T) {
	Run(TestConfig{
		Name: "kustomize",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("kustomize version")
			return nil
		},
	})(t)
}
