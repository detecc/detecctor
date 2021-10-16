package api

type (
	CommandBuilder struct {
		buildActions []handler
	}

	handler func(cmd *Command)
)

//NewCommandBuilder - constructor
func NewCommandBuilder() *CommandBuilder {
	return &CommandBuilder{}
}

//WithName sets name of the command
func (b *CommandBuilder) WithName(value string) *CommandBuilder {
	b.buildActions = append(b.buildActions, func(cmd *Command) {
		cmd.Name = value
	})
	return b
}

//WithArgs sets arguments of the command
func (b *CommandBuilder) WithArgs(args []string) *CommandBuilder {
	b.buildActions = append(b.buildActions, func(cmd *Command) {
		cmd.Args = args
	})
	return b
}

//FromChat sets chatId
func (b *CommandBuilder) FromChat(chatId string) *CommandBuilder {
	b.buildActions = append(b.buildActions, func(cmd *Command) {
		cmd.ChatId = chatId
	})
	return b
}

//Build builds the Command object
func (b *CommandBuilder) Build() Command {
	emp := Command{
		Name:   "",
		Args:   []string{},
		ChatId: "0",
	}

	for _, a := range b.buildActions {
		a(&emp)
	}
	return emp
}
