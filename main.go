package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/juju/errgo"
	"github.com/spf13/cobra"

	"./service"
)

var (
	projectVersion = "dev"
	projectBuild   = "dev"

	maskAny = errgo.MaskFunc(errgo.Any)
)

var (
	cmdMain = &cobra.Command{
		Use: "vault-monkey",
		Run: showUsage,
	}
	globalFlags struct {
		debug   bool
		verbose bool
		service.VaultServiceConfig
	}
)

func init() {
	cmdMain.PersistentFlags().BoolVarP(&globalFlags.debug, "debug", "D", false, "Print debug output")
	cmdMain.PersistentFlags().BoolVarP(&globalFlags.verbose, "verbose", "v", false, "Print verbose output")
	cmdMain.PersistentFlags().StringVar(&globalFlags.VaultAddr, "vault-addr", "", "URL of the vault (defaults to VAULT_ADDR environment variable)")
	cmdMain.PersistentFlags().StringVar(&globalFlags.TokenPath, "token-path", "", "Path of a file containing your vault token (token defaults to VAULT_TOKEN environment variable)")
}

func main() {
	cmdMain.Execute()
}

func showUsage(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

func Exitf(format string, args ...interface{}) {
	if !strings.HasSuffix(format, "\n") {
		format = format + "\n"
	}
	fmt.Printf(format, args...)
	fmt.Println()
	os.Exit(1)
}

func Verbosef(format string, args ...interface{}) {
	if globalFlags.verbose {
		fmt.Printf(format, args...)
	}
}

func assert(err error) {
	if err != nil {
		Exitf("Assertion failed: %#v", err)
	}
}

func assertArgIsSet(arg, argKey string) {
	if arg == "" {
		Exitf("%s must be set\n", argKey)
	}
}
