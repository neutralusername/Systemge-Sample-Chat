package appChat

import "github.com/neutralusername/Systemge/Node"

func (app *App) OnStart(node *Node.Node) error {
	app.rooms = map[string]*room{}
	app.chatters = map[string]*chatter{}
	return nil
}
