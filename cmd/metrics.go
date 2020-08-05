package cmd

import (
	"encoding/json"
	"fmt"
	"grpctl/transport"
	"log"
	"strings"

	"github.com/spf13/cobra"
	zpb "google.golang.org/grpc/channelz/grpc_channelz_v1"
)

type channelzMetrics struct {
	Channels       []*zpb.Channel    `json:"channels,omitempty"`
	Subchannels    []*zpb.Subchannel `json:"subchannels,omitempty"`
	Servers        []*zpb.Server     `json:"servers,omitempty"`
	ServerSockets  []*zpb.Socket     `json:"server_sockets,omitempty"`
	ChannelSockets []*zpb.Socket     `json:"channel_sockets,omitempty"`
}

func addChannels(m *channelzMetrics) {
	m.Channels = transport.Channels()
}

func addSubchannels(m *channelzMetrics) {
	m.Subchannels = transport.Subchannels()
}

func addServers(m *channelzMetrics) {
	m.Servers = transport.Servers()
}

func addServerSockets(m *channelzMetrics) {
	m.ServerSockets = transport.ServerSockets()
}

func addChannelSockets(m *channelzMetrics) {
	m.ChannelSockets = transport.ChannelSockets()
}

var metricsAll []string = []string{"channels", "subchannels", "servers", "server-sockets", "channel-sockets"}
var metricsArgActions map[string]func(m *channelzMetrics) = map[string]func(m *channelzMetrics){
	"channels":        addChannels,
	"subchannels":     addSubchannels,
	"servers":         addServers,
	"server-sockets":  addServerSockets,
	"channel-sockets": addChannelSockets,
}

func metricsCommandRunWithError(cmd *cobra.Command, args []string) error {
	if !transport.IsConnected() {
		return errNoWatch
	}
	var m channelzMetrics

	if len(args) == 0 {
		args = metricsAll
	}
	for _, arg := range args {
		metricsArgActions[arg](&m)
	}

	r, err := json.MarshalIndent(&m, "", "\t")
	if err != nil {
		log.Fatalf("failed to marshal metrics: %v", err)
	}
	fmt.Println(string(r))
	return nil
}

var metricsCmd = &cobra.Command{
	Use:   "metrics",
	Short: "Pull metrics from the service being watched",
	Long: fmt.Sprintf(`Pull metrics from the service being watched. Available metric categories are:

	%v`, strings.Join(metricsAll, ", ")),
	Args: func(cmd *cobra.Command, args []string) error {
		for _, arg := range args {
			if _, ok := metricsArgActions[arg]; !ok {
				return fmt.Errorf("Unknown metrics category: %v", arg)
			}
		}
		return nil
	},
	RunE: metricsCommandRunWithError,
}

func init() {
	rootCmd.AddCommand(metricsCmd)
}
