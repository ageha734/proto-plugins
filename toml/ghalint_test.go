package main

import (
	"testing"
)

func TestGhalint(t *testing.T) {
	Run(TestConfig{
		Name: "ghalint",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("ghalint --version")
			return nil
		},
	})(t)
}
