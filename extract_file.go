package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"git.pulcy.com/pulcy/vault-monkey/service"
)

var (
	cmdExtractFile = &cobra.Command{
		Use:     "file",
		Short:   "Extract a secret into a file.",
		Example: "--target=<file-path> <secret-path>[#<secret-field>]",
		Run:     cmdExtractFileRun,
	}
)

func init() {
	cmdExtract.AddCommand(cmdExtractFile)
}

func cmdExtractFileRun(cmd *cobra.Command, args []string) {
	// Check arguments
	if len(args) != 1 {
		Exitf("Provide exactly one argument: <path>[#field]")
	}

	// Parse arguments
	secretPath, secretField, err := parseSecretPath(args[0])
	if err != nil {
		Exitf(err.Error())
	}

	// Login
	vs, err := serverLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	// Create env file
	secret := service.FileSecret{
		SecretPath:  secretPath,
		SecretField: secretField,
	}
	if err := vs.CreateSecretFile(extractFlags.targetFilePath, secret); err != nil {
		Exitf("Secret extraction failed: %v", err)
	}
}

func parseSecretPath(arg string) (string, string, error) {
	pf := strings.Split(arg, "#")
	switch len(pf) {
	case 1:
		return arg, "value", nil
	case 2:
		return pf[0], pf[1], nil
	default:
		return "", "", maskAny(fmt.Errorf("expected '<path>[#field]', got '%s'", arg))
	}
}
