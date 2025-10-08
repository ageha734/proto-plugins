package main

import (
	"testing"
)

func TestKubens(t *testing.T) {
	Run(TestConfig{
		Name: "kubens",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("kubens --version")
			return nil
		},
	})(t)
}
