package main

import (
	"testing"
)

func TestTilt(t *testing.T) {
	Run(TestConfig{
		Name: "tilt",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("tilt version")
			return nil
		},
	})(t)
}
