package transport

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	zpb "google.golang.org/grpc/channelz/grpc_channelz_v1"
	"google.golang.org/grpc/connectivity"
)

var conn *grpc.ClientConn
var channelzClient zpb.ChannelzClient

// Connect connects to the service at address and creates stubs
func Connect(address string) {
	var err error
	conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	channelzClient = zpb.NewChannelzClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	var state connectivity.State = conn.GetState()
	for state != connectivity.Ready {
		conn.WaitForStateChange(ctx, state)
		if ctx.Err() != nil {
			log.Fatalf("failed to establish connection to address: %v", address)
		}
		state = conn.GetState()
	}
}

// IsConnected checks if the connection is ready to use
func IsConnected() bool {
	if conn == nil {
		return false
	}
	return conn.GetState() == connectivity.Ready
}

// Channels returns all available channels
func Channels() []*zpb.Channel {
	channels, err := channelzClient.GetTopChannels(context.Background(), &zpb.GetTopChannelsRequest{})
	if err != nil {
		log.Fatalf("failed to fetch top channels: %v", err)
	}
	return channels.Channel
}

// Subchannel returns the queried subchannel
func Subchannel(subchannelID int64) *zpb.Subchannel {
	subchannel, err := channelzClient.GetSubchannel(context.Background(), &zpb.GetSubchannelRequest{SubchannelId: subchannelID})
	if err != nil {
		log.Fatalf("failed to fetch subchannel (id=%v): %v", subchannelID, err)
	}
	return subchannel.Subchannel
}

// Subchannels traverses all channels and fetches all subchannels
func Subchannels() []*zpb.Subchannel {
	var s []*zpb.Subchannel
	for _, channel := range Channels() {
		for _, subchannelRef := range channel.SubchannelRef {
			s = append(s, Subchannel(subchannelRef.SubchannelId))
		}
	}
	return s
}

// Servers returns all available servers
func Servers() []*zpb.Server {
	servers, err := channelzClient.GetServers(context.Background(), &zpb.GetServersRequest{})
	if err != nil {
		log.Fatalf("failed to fetch servers: %v", err)
	}
	return servers.Server
}

// Socket returns a socket
func Socket(socketID int64) *zpb.Socket {
	socket, err := channelzClient.GetSocket(context.Background(), &zpb.GetSocketRequest{SocketId: socketID})
	if err != nil {
		log.Fatalf("failed to fetch socket (id=%v): %v", socketID, err)
	}
	return socket.Socket
}

// ServerSockets returns all sockets for servers
func ServerSockets() []*zpb.Socket {
	var s []*zpb.Socket
	for _, server := range Servers() {
		serverSocketResp, err := channelzClient.GetServerSockets(
			context.Background(),
			&zpb.GetServerSocketsRequest{ServerId: server.Ref.ServerId},
		)
		if err != nil {
			log.Fatalf("failed to fetch server sockets (id=%v): %v", server.Ref.ServerId, err)
		}
		for _, socketRef := range serverSocketResp.SocketRef {
			s = append(s, Socket(socketRef.SocketId))
		}
	}
	return s
}

// ChannelSockets returns all sockets for clients (channels)
func ChannelSockets() []*zpb.Socket {
	var s []*zpb.Socket
	for _, subchannel := range Subchannels() {
		for _, socketRef := range subchannel.SocketRef {
			s = append(s, Socket(socketRef.SocketId))
		}
	}
	return s
}
