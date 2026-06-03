package main

import (
	"testing"
)

func TestArgocd(t *testing.T) {
	Run(TestConfig{
		Name: "argocd",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("argocd version --client")
			return nil
		},
	})(t)
}
