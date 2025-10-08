package main

import (
	"testing"
)

func TestTrivy(t *testing.T) {
	Run(TestConfig{
		Name: "trivy",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("trivy --version")
			return nil
		},
	})(t)
}
