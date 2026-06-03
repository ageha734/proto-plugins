package main

import (
	"testing"
)

func TestTerraformer(t *testing.T) {
	Run(TestConfig{
		Name: "terraformer",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("terraformer version")
			return nil
		},
	})(t)
}
