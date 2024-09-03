package main

import (
	"SystemgeSampleChat/appChat"
	"SystemgeSampleChat/appWebsocketHttp"
	"SystemgeSampleChat/topics"
	"time"

	"github.com/neutralusername/Systemge/BrokerResolver"
	"github.com/neutralusername/Systemge/BrokerServer"
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Error"
)

const LOGGER_PATH = "logs.log"

func main() {
	if err := Dashboard.NewServer("DasbhboardServer",
		&Config.DashboardServer{
			HTTPServerConfig: &Config.HTTPServer{
				TcpServerConfig: &Config.TcpServer{
					Port: 8081,
				},
			},
			WebsocketServerConfig: &Config.WebsocketServer{
				Pattern:                 "/ws",
				ClientWatchdogTimeoutMs: 1000 * 60,
				TcpServerConfig: &Config.TcpServer{
					Port: 8444,
				},
			},
			SystemgeServerConfig: &Config.SystemgeServer{
				ListenerConfig: &Config.TcpListener{
					TcpServerConfig: &Config.TcpServer{
						Port: 60000,
					},
				},
				ConnectionConfig: &Config.TcpConnection{},
			},
			HeapUpdateIntervalMs:      1000,
			GoroutineUpdateIntervalMs: 1000,
			StatusUpdateIntervalMs:    1000,
			MetricsUpdateIntervalMs:   1000,
			MaxChartEntries:           100,
		},
	).Start(); err != nil {
		panic(Error.New("Dashboard server failed to start", err))
	}
	if err := BrokerResolver.New("brokerResolver",
		&Config.MessageBrokerResolver{
			SystemgeServerConfig: &Config.SystemgeServer{
				ListenerConfig: &Config.TcpListener{
					TcpServerConfig: &Config.TcpServer{
						Port: 60001,
					},
				},
				ConnectionConfig: &Config.TcpConnection{},
			},
			DashboardClientConfig: &Config.DashboardClient{
				ConnectionConfig: &Config.TcpConnection{},
				ClientConfig: &Config.TcpClient{
					Address: "localhost:60000",
				},
			},
			AsyncTopicClientConfigs: map[string]*Config.TcpClient{
				topics.PROPAGATE_MESSAGE: {
					Address: "localhost:60002",
				},
				topics.ADD_MESSAGE: {
					Address: "localhost:60002",
				},
			},
			SyncTopicClientConfigs: map[string]*Config.TcpClient{
				topics.JOIN: {
					Address: "localhost:60002",
				},
				topics.LEAVE: {
					Address: "localhost:60002",
				},
			},
		},
	).Start(); err != nil {
		panic(Error.New("MessageBroker resolver failed to start", err))
	}

	if err := BrokerServer.New("brokerServer",
		&Config.MessageBrokerServer{
			SystemgeServerConfig: &Config.SystemgeServer{
				ListenerConfig: &Config.TcpListener{
					TcpServerConfig: &Config.TcpServer{
						Port: 60002,
					},
				},
				ConnectionConfig: &Config.TcpConnection{},
			},
			AsyncTopics: []string{topics.PROPAGATE_MESSAGE, topics.ADD_MESSAGE},
			SyncTopics:  []string{topics.JOIN, topics.LEAVE},
			DashboardClientConfig: &Config.DashboardClient{
				ConnectionConfig: &Config.TcpConnection{},
				ClientConfig: &Config.TcpClient{
					Address: "localhost:60000",
				},
			},
		},
	).Start(); err != nil {
		panic(Error.New("MessageBroker server failed to start", err))
	}
	appWebsocketHttp.New()

	appChat.New()
	<-make(chan time.Time)
}
