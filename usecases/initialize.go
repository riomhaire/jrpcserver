package usecases

import "github.com/riomhaire/jrpcserver/model"

var Commands []model.JRPCCommand

func InitializeCommands() []model.JRPCCommand {
	commands := make([]model.JRPCCommand, 0)

	commands = append(commands, model.JRPCCommand{"test.ping", PingCommand})
	commands = append(commands, model.JRPCCommand{"test.pong", PongCommand})
	commands = append(commands, model.JRPCCommand{"test.echo", EchoCommand})
	commands = append(commands, model.JRPCCommand{"system.commands", ListCommandsCommand})
	Commands = commands // needed for list
	return commands
}
