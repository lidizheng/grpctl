package cmd

import (
	"fmt"
	"grpctl/store"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "grpctl",
	Short: "grpctl is an gRPC service admin CLI",
}
var verboseFlag bool

// Execute executes the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	// After the verbose is set, we can use Debugf to log
	if verboseFlag {
		store.SetVerbose()
	}

	maybeStartWatching()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolVarP(&verboseFlag, "verbose", "v", false, "Print verbose information for debugging")
}
