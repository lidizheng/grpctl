package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"grpctl/transport"
	"time"

	"github.com/spf13/cobra"
)

type buildInfo struct {
	BuildTime         time.Time `json:"build_time,string"`
	CompilerInfo      string    `json:"compiler_info"`
	SourceCodeVersion string    `json:"source_code_version"`
}

var buildInfoCmd = &cobra.Command{
	Use:   "build",
	Short: "Print the build information",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !transport.IsConnected() {
			return errNoWatch
		}

		r, err := json.MarshalIndent(defaultBuildInfo, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(r))
		return nil
	},
}

type xdsDescription struct {
	Listeners []string `json:"listeners"`
	Routes    []string `json:"routes"`
	Clusters  []string `json:"clusters"`
}

type deploymentInfo struct {
	DeploymentTime     time.Time      `json:"deployment_time,string"`
	XdsSelfDescription xdsDescription `json:"xds_self_description"`
	Certificates       []string       `json:"certificates"`
}

var deploymentInfoCmd = &cobra.Command{
	Use:   "deployment",
	Short: "Print the deployment information",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !transport.IsConnected() {
			return errNoWatch
		}

		r, err := json.MarshalIndent(defaultDeploymentInfo, "", "\t")
		if err != nil {
			return err
		}
		fmt.Println(string(r))
		return nil
	},
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Print information about the service being watched.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !transport.IsConnected() {
			return errNoWatch
		}
		return errors.New("please specify what information to get")
	},
}

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check the health of the peer service.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !transport.IsConnected() {
			return errNoWatch
		}
		// Todo
		if len(args) == 0 {
			fmt.Println("OK")
		} else {
			for _, arg := range args {
				fmt.Printf("%v: OK\n", arg)
			}
		}
		return nil
	},
}

var defaultBuildInfo = buildInfo{
	BuildTime:         time.Date(2020, time.April, 1, 8, 0, 0, 0, time.Local),
	CompilerInfo:      "gcc-4.X.Y-crosstool-v18-llvm-grtev4-k8",
	SourceCodeVersion: "9106c3fff5236fd664a8de183f1c27682c66b823",
}

var defaultDeploymentInfo = deploymentInfo{
	DeploymentTime: time.Date(2020, time.April, 1, 9, 0, 0, 0, time.Local),
	XdsSelfDescription: xdsDescription{
		Listeners: []string{"TRAFFICDIRECTOR_INTERCEPTION_LISTENER"},
		Routes:    []string{"URL_MAP/830293263384_lidiz-td-url-map"},
		Clusters:  []string{"cloud-internal-istio:cloud_mp_830293263384_7265588390596584651"},
	},
	Certificates: nil,
}

func init() {
	infoCmd.AddCommand(buildInfoCmd)
	infoCmd.AddCommand(deploymentInfoCmd)
	infoCmd.AddCommand(healthCmd)
	rootCmd.AddCommand(infoCmd)
}
