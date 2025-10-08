package main

import (
	"testing"
)

func TestLefthook(t *testing.T) {
	Run(TestConfig{
		Name: "lefthook",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("lefthook version")
			return nil
		},
	})(t)
}
