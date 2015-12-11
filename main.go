package main

import (
	"fmt"
	"os"

	"github.com/juju/errgo"
	"github.com/spf13/cobra"
)

var (
	projectVersion = "dev"
	projectBuild   = "dev"

	maskAny = errgo.MaskFunc(errgo.Any)
)

var (
	cmdMain = &cobra.Command{
		Use: "secrext",
		Run: showUsage,
	}
	globalFlags struct {
		debug          bool
		verbose        bool
		secretPath     string
		secretField    string
		targetFilePath string
	}
)

func init() {
	cmdMain.PersistentFlags().BoolVarP(&globalFlags.debug, "debug", "D", false, "Print debug output")
	cmdMain.PersistentFlags().BoolVarP(&globalFlags.verbose, "verbose", "v", false, "Print verbose output")
	cmdMain.PersistentFlags().StringVar(&globalFlags.secretPath, "path", "", "Path of the secret in vault")
	cmdMain.PersistentFlags().StringVar(&globalFlags.secretField, "field", "", "Field within the secret in vault")
	cmdMain.PersistentFlags().StringVar(&globalFlags.targetFilePath, "target", "", "Path of target file")
}

func main() {
	cmdMain.AddCommand(envCmd)
	cmdMain.AddCommand(fileCmd)

	cmdMain.Execute()
}

func showUsage(cmd *cobra.Command, args []string) {
	cmd.Usage()
}

func Exitf(format string, args ...interface{}) {
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
