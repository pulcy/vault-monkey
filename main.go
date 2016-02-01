package main

import (
	"fmt"
	"os"
	"strings"

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
		Use: "vault-monkey",
		Run: showUsage,
	}
	globalFlags struct {
		debug   bool
		verbose bool
	}
)

func init() {
	cmdMain.PersistentFlags().BoolVarP(&globalFlags.debug, "debug", "D", false, "Print debug output")
	cmdMain.PersistentFlags().BoolVarP(&globalFlags.verbose, "verbose", "v", false, "Print verbose output")
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
