package main

import (
	"testing"
)

func TestKubeconform(t *testing.T) {
	Run(TestConfig{
		Name: "kubeconform",
		AfterInstall: func(t *testing.T, shell *Shell) error {
			shell.Exec("kubeconform -v")
			return nil
		},
	})(t)
}
