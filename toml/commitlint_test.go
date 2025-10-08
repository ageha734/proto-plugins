package main

import (
	"testing"
)

func TestCommitlint(t *testing.T) {
	Run(TestConfig{
		Name: "commitlint",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("commitlint --version")
			return nil
		},
	})(t)
}
