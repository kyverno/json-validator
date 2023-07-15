package main

import "github.com/kyverno/json-validator/backend/pkg/cmd"

func main() {
	rootCmd := cmd.NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
