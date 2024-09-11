package main

import (
	"SystemgeSampleChat/appChat"
	"SystemgeSampleChat/appWebsocketHttp"
	"SystemgeSampleChat/topics"
	"time"

	"github.com/neutralusername/Systemge/BrokerResolver"
	"github.com/neutralusername/Systemge/BrokerServer"
	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/DashboardClientCustomService"
	"github.com/neutralusername/Systemge/DashboardServer"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
)

const LOGGER_PATH = "logs.log"

func main() {
	Helpers.StartPprof()
	if err := DashboardServer.New("DasbhboardServer",
		&Config.DashboardServer{
			HTTPServerConfig: &Config.HTTPServer{
				TcpServerConfig: &Config.TcpServer{
					Port: 8081,
				},
			},
			WebsocketServerConfig: &Config.WebsocketServer{
				Pattern:                 "/ws",
				ClientWatchdogTimeoutMs: 1000 * 60 * 10,
				TcpServerConfig: &Config.TcpServer{
					Port: 8444,
				},
			},
			SystemgeServerConfig: &Config.SystemgeServer{
				TcpSystemgeListenerConfig: &Config.TcpSystemgeListener{
					TcpServerConfig: &Config.TcpServer{
						Port: 60000,
					},
				},
				TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{
					HeartbeatIntervalMs: 1000,
				},
			},
			DashboardMetrics:           true,
			DashboardCommands:          true,
			DashboardSystemgeCommands:  true,
			DashboardHttpCommands:      true,
			DashboardWebsocketCommands: true,
			HeapUpdateIntervalMs:       1000,
			GoroutineUpdateIntervalMs:  1000,
			StatusUpdateIntervalMs:     1000,
			MetricsUpdateIntervalMs:    1000,
			MaxChartEntries:            100,
		},
		nil, nil,
	).Start(); err != nil {
		panic(Error.New("Dashboard server failed to start", err))
	}

	brokerResolver := BrokerResolver.New("brokerResolver",
		&Config.MessageBrokerResolver{
			SystemgeServerConfig: &Config.SystemgeServer{
				TcpSystemgeListenerConfig: &Config.TcpSystemgeListener{
					TcpServerConfig: &Config.TcpServer{
						Port: 60001,
					},
				},
				TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{
					HeartbeatIntervalMs: 1000,
				},
			},
			AsyncTopicClientConfigs: map[string]*Config.TcpClient{
				topics.PROPAGATE_MESSAGE: {
					Address: "[::1]:60002",
				},
				topics.ADD_MESSAGE: {
					Address: "[::1]:60002",
				},
			},
			SyncTopicClientConfigs: map[string]*Config.TcpClient{
				topics.JOIN: {
					Address: "[::1]:60002",
				},
				topics.LEAVE: {
					Address: "[::1]:60002",
				},
			},
		},
		nil, nil,
	)
	if err := DashboardClientCustomService.New("brokerResolver", &Config.DashboardClient{
		TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{
			HeartbeatIntervalMs: 1000,
		},
		TcpClientConfig: &Config.TcpClient{
			Address: "[::1]:60000",
		},
	}, brokerResolver, nil).Start(); err != nil {
		panic(Error.New("Dashboard client failed to start", err))
	}

	brokerServer := BrokerServer.New("brokerServer",
		&Config.MessageBrokerServer{
			SystemgeServerConfig: &Config.SystemgeServer{
				TcpSystemgeListenerConfig: &Config.TcpSystemgeListener{
					TcpServerConfig: &Config.TcpServer{
						Port: 60002,
					},
				},
				TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{
					HeartbeatIntervalMs: 1000,
				},
			},
			AsyncTopics: []string{topics.PROPAGATE_MESSAGE, topics.ADD_MESSAGE},
			SyncTopics:  []string{topics.JOIN, topics.LEAVE},
		},
		nil, nil,
	)

	if err := DashboardClientCustomService.New("brokerServer", &Config.DashboardClient{
		TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{
			HeartbeatIntervalMs: 1000,
		},
		TcpClientConfig: &Config.TcpClient{
			Address: "[::1]:60000",
		},
	}, brokerServer, nil).Start(); err != nil {
		panic(Error.New("Dashboard client failed to start", err))
	}

	appWebsocketHttp.New()
	appChat.New()
	<-make(chan time.Time)
}
