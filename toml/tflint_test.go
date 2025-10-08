package main

import (
	"testing"
)

func TestTflint(t *testing.T) {
	Run(TestConfig{
		Name: "tflint",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("tflint --version")
			return nil
		},
	})(t)
}
