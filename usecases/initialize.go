package usecases

import (
	"github.com/riomhaire/jrpcserver/model"
	"github.com/riomhaire/jrpcserver/usecases/defaultcommand"
)

var Commands []model.JRPCCommand

func InitializeCommands() []model.JRPCCommand {
	commands := make([]model.JRPCCommand, 0)

	commands = append(commands, model.JRPCCommand{"test.ping", defaultcommand.PingCommand, false})
	commands = append(commands, model.JRPCCommand{"test.pong", defaultcommand.PongCommand, false})
	commands = append(commands, model.JRPCCommand{"test.echo", defaultcommand.EchoCommand, true})
	commands = append(commands, model.JRPCCommand{"system.commands", defaultcommand.ListCommandsCommand, false})
	Commands = commands // needed for list
	return commands
}
