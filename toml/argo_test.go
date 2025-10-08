package main

import (
	"testing"
)

func TestArgo(t *testing.T) {
	Run(TestConfig{
		Name: "argo",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("argo version")
			return nil
		},
	})(t)
}
