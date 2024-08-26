package main

import (
	"SystemgeSampleChat/appChat"
	"SystemgeSampleChat/appWebsocketHttp"
	"SystemgeSampleChat/topics"
	"time"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/MessageBroker"
)

const LOGGER_PATH = "logs.log"

func main() {
	if Dashboard.NewServer(&Config.DashboardServer{
		HTTPServerConfig: &Config.HTTPServer{
			TcpListenerConfig: &Config.TcpListener{
				Port: 8081,
			},
		},
		WebsocketServerConfig: &Config.WebsocketServer{
			Pattern:                 "/ws",
			ClientWatchdogTimeoutMs: 1000 * 60,
			TcpListenerConfig: &Config.TcpListener{
				Port: 8444,
			},
		},
		SystemgeServerConfig: &Config.SystemgeServer{
			Name: "dashboardServer",
			ListenerConfig: &Config.SystemgeListener{
				TcpListenerConfig: &Config.TcpListener{
					Port: 60000,
				},
			},
			ConnectionConfig: &Config.SystemgeConnection{},
		},
		HeapUpdateIntervalMs:      1000,
		GoroutineUpdateIntervalMs: 1000,
		StatusUpdateIntervalMs:    1000,
		MetricsUpdateIntervalMs:   1000,
	}).Start() != nil {
		panic("Dashboard server failed to start")
	}
	if MessageBroker.NewMessageBrokerServer(&Config.MessageBrokerServer{
		SystemgeServerConfig: &Config.SystemgeServer{
			Name: "messageBrokerServer",
			ListenerConfig: &Config.SystemgeListener{
				TcpListenerConfig: &Config.TcpListener{
					Port: 60001,
				},
			},
			ConnectionConfig: &Config.SystemgeConnection{},
		},
		AsyncTopics: []string{topics.PROPAGATE_MESSAGE, topics.ADD_MESSAGE},
		SyncTopics:  []string{topics.JOIN, topics.LEAVE},
		DashboardClientConfig: &Config.DashboardClient{
			Name:             "messageBrokerServer",
			ConnectionConfig: &Config.SystemgeConnection{},
			EndpointConfig: &Config.TcpEndpoint{
				Address: "localhost:60000",
			},
		},
	}).Start() != nil {
		panic("MessageBroker server failed to start")
	}
	appWebsocketHttp.New()

	appChat.New()
	<-make(chan time.Time)
}
