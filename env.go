package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var (
	envCmd = &cobra.Command{
		Use:   "env",
		Short: "Extract a secret into an environment file.",
		Run:   runEnv,
	}
	envFlags struct {
		KeyName string
	}
)

func init() {
	envCmd.Flags().StringVar(&envFlags.KeyName, "key", "", "key of environment variable")
}

func runEnv(cmd *cobra.Command, args []string) {
	if envFlags.KeyName == "" {
		Exitf("key not set\n")
	}
	// get secret
	secret, err := extractSecret(globalFlags.secretPath, globalFlags.secretField)
	if err != nil {
		Exitf("Failed to extract secret: %v\n", err)
	}

	// format environment file
	content := fmt.Sprintf("%s=%s\n", envFlags.KeyName, strconv.Quote(secret))

	// save environment file
	if err := saveTargetFile(globalFlags.targetFilePath, content); err != nil {
		Exitf("Failed to created %s: %v\n", globalFlags.targetFilePath, err)
	}
}
