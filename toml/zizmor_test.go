package main

import (
	"testing"
)

func TestZizmor(t *testing.T) {
	Run(TestConfig{
		Name: "zizmor",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("zizmor --version")
			return nil
		},
	})(t)
}
