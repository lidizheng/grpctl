package cmd

import "github.com/spf13/cobra"

var actionCmd = &cobra.Command{
	Use:   "action",
	Short: "Manually operate the target gRPC service - todo",
	Run:   func(cmd *cobra.Command, args []string) {},
}

func init() {
	rootCmd.AddCommand(actionCmd)
}
