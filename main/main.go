package main

import (
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Utilities"
	"SystemgeSampleChat/appChat"
	"SystemgeSampleChat/appWebsocketHTTP"
)

const RESOLVER_ADDRESS = "127.0.0.1:60000"
const RESOLVER_NAME_INDICATION = "127.0.0.1"
const RESOLVER_TLS_CERT_PATH = "MyCertificate.crt"
const WEBSOCKET_PORT = ":8443"
const HTTP_PORT = ":8080"

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	//resolver and brokers are placed outside, because on startup they need to be started first.
	//and if they are stopped first on stop, which is inherent to multi-module behaviour, the other modules will not be able to communicate with them during their stop/disconnect routine.
	//this demonstrates why multi modules are not always the best solution.
	//the alternative is to start them either manually, like here in the main function, or to start them in separate terminal windows as separate processes/programs.
	err := Module.NewResolverFromConfig("resolver.systemge", ERROR_LOG_FILE_PATH).Start()
	if err != nil {
		panic(err)
	}
	err = Module.NewBrokerFromConfig("brokerChat.systemge", ERROR_LOG_FILE_PATH).Start()
	if err != nil {
		panic(err)
	}
	err = Module.NewBrokerFromConfig("brokerWebsocket.systemge", ERROR_LOG_FILE_PATH).Start()
	if err != nil {
		panic(err)
	}

	nodeChat := Module.NewNode(&Node.Config{
		Name:                   "nodeApp",
		ResolverAddress:        RESOLVER_ADDRESS,
		ResolverNameIndication: RESOLVER_NAME_INDICATION,
		ResolverTLSCert:        Utilities.GetFileContent(RESOLVER_TLS_CERT_PATH),
		LoggerPath:             ERROR_LOG_FILE_PATH,
	}, appChat.New(), nil, nil)
	appWebsocketHTTP := appWebsocketHTTP.New()
	nodeWebsocket := Module.NewNode(&Node.Config{
		Name:                   "nodeWebsocketHTTP",
		ResolverAddress:        RESOLVER_ADDRESS,
		ResolverNameIndication: RESOLVER_NAME_INDICATION,
		ResolverTLSCert:        Utilities.GetFileContent(RESOLVER_TLS_CERT_PATH),
		WebsocketPattern:       "/ws",
		WebsocketPort:          WEBSOCKET_PORT,
		HTTPPort:               HTTP_PORT,
		LoggerPath:             ERROR_LOG_FILE_PATH,
	}, appWebsocketHTTP, appWebsocketHTTP, appWebsocketHTTP)
	Module.StartCommandLineInterface(Module.NewMultiModule(
		nodeWebsocket,
		nodeChat,
	))
}
