package main

import (
	"testing"
)

func TestKubectl(t *testing.T) {
	Run(TestConfig{
		Name: "kubectl",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("kubectl version")
			return nil
		},
	})(t)
}
