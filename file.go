package main

import (
	"github.com/spf13/cobra"
)

var (
	fileCmd = &cobra.Command{
		Use:   "file",
		Short: "Extract a secret into a file.",
		Run:   runFile,
	}
)

func runFile(cmd *cobra.Command, args []string) {
	// get secret
	secret, err := extractSecret(globalFlags.secretPath, globalFlags.secretField)
	if err != nil {
		Exitf("Failed to extract secret: %v\n", err)
	}

	// save environment file
	if err := saveTargetFile(globalFlags.targetFilePath, secret); err != nil {
		Exitf("Failed to created %s: %v\n", globalFlags.targetFilePath, err)
	}
}
