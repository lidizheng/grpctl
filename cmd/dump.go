package cmd

import (
	"github.com/spf13/cobra"
)

var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Print gRPC configurations - todo",
	Run:  func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(dumpCmd)
}
