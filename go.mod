module SystemgeSampleChat

go 1.22.3

replace Systemge => ../Systemge

require (
	github.com/gorilla/websocket v1.5.3
	github.com/neutralusername/Systemge v0.0.0-20240727153215-3150bb25c175
)

require golang.org/x/oauth2 v0.21.0 // indirect
