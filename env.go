package main

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"./service"
)

var (
	cmdExtractEnv = &cobra.Command{
		Use:     "env",
		Short:   "Extract a secret into an environment file.",
		Example: "--target=<file-path> <environment-key>=<secret-path>[#<secret-field>]...",
		Run:     cmdExtractEnvRun,
	}
)

func init() {
	cmdExtract.AddCommand(cmdExtractEnv)
}

func cmdExtractEnvRun(cmd *cobra.Command, args []string) {
	// Check arguments
	if len(args) == 0 {
		Exitf("Private at least one argument: <key>=<path>[#field]")
	}

	// Parse arguments
	secrets := []service.EnvSecret{}
	for _, arg := range args {
		secret, err := parseEnvSecret(arg)
		if err != nil {
			Exitf(err.Error())
		}
		secrets = append(secrets, secret)
	}

	// Login
	vs, err := serverLogin()
	if err != nil {
		Exitf("Login failed: %v", err)
	}

	// Create env file
	if err := vs.CreateEnvironmentFile(extractFlags.targetFilePath, secrets); err != nil {
		Exitf("Secret extraction failed: %v", err)
	}
}

func parseEnvSecret(arg string) (service.EnvSecret, error) {
	kv := strings.Split(arg, "=")
	if len(kv) != 2 {
		return service.EnvSecret{}, maskAny(fmt.Errorf("expected '<key>=<path>[#field]', got '%s'", arg))
	}
	envKey := kv[0]
	secretPath, secretField, err := parseSecretPath(kv[1])
	if err != nil {
		return service.EnvSecret{}, maskAny(err)
	}
	return service.EnvSecret{
		SecretPath:     secretPath,
		SecretField:    secretField,
		EnvironmentKey: envKey,
	}, nil
}
