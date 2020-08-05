package cmd

import (
	"fmt"
	"grpctl/store"
	"grpctl/transport"
	"os"

	"github.com/spf13/cobra"
)

var watchServiceAddr string

const watchEnvName string = "GRPCTL_WATCH"

var errNoWatch error = fmt.Errorf("please specify watch a target service through --watch argument or environment variable %v", watchEnvName)

func maybeStartWatching() {
	var target string
	if watchServiceAddr != "" {
		target = watchServiceAddr
		store.Debugf("watch address supplied by commandline argument: %v", target)
	}

	if target == "" && os.Getenv(watchEnvName) != "" {
		target = os.Getenv(watchEnvName)
		store.Debugf("watch address supplied by environment variable: %v", target)
	}

	if target == "" {
		target = store.LoadTarget()
		if target != "" {
			store.Debugf("watch address supplied by config file: %v", target)
		}
	}

	if target != "" {
		transport.Connect(target)
	}
}

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Save the target address.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		target := args[0]
		if target != "" {
			transport.Connect(target)
			return store.SaveTarget(target)
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&watchServiceAddr, "watch", "w", "", "Address of the target service")
	rootCmd.AddCommand(watchCmd)
}
