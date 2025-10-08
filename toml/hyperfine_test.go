package main

import (
	"testing"
)

func TestHyperfine(t *testing.T) {
	Run(TestConfig{
		Name: "hyperfine",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("hyperfine --version")
			return nil
		},
	})(t)
}
