package main

import (
	"testing"
)

func TestTask(t *testing.T) {
	Run(TestConfig{
		Name: "task",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("task --version")
			return nil
		},
	})(t)
}
