package main

import (
	"testing"
)

func TestHadolint(t *testing.T) {
	Run(TestConfig{
		Name: "hadolint",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("hadolint --version")
			return nil
		},
	})(t)
}
