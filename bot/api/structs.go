package api

// constants for the Reply type
const (
	TypeMessage = 0
	TypePhoto   = 1
	TypeAudio   = 2
)

type (
	// Reply is a struct used to parse results to the ReplyChannel in Bot.
	Reply struct {
		// Each reply must contain a ChatId - a chat to reply to.
		ChatId string
		// The ReplyType must be a constant defined in the package.
		ReplyType int
		// Content must be cast after determining the type to send to Bot.
		Content interface{}
	}

	// Command consists of a Name and Args.
	// Example of a command: "/get_status serviceNode1 serviceNode2".
	// The command name is "/get_status", the arguments are ["serviceNode1", "serviceNode2"].
	Command struct {
		Name string
		// Args contains the arguments extracted from the Message sent through Telegram.
		Args   []string
		ChatId string
	}
)
