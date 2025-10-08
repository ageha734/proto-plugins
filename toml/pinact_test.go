package main

import (
	"testing"
)

func TestPinact(t *testing.T) {
	Run(TestConfig{
		Name: "pinact",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("pinact --version")
			return nil
		},
	})(t)
}
