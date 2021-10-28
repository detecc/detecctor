package api

type (
	ReplyBuilder struct {
		buildActions []replyHandler
	}

	replyHandler func(r *Reply)
)

//NewReplyBuilder - constructor
func NewReplyBuilder() *ReplyBuilder {
	return &ReplyBuilder{}
}

//TypeMessage sets the type to TypeMessage
func (b *ReplyBuilder) TypeMessage() *ReplyBuilder {
	b.buildActions = append(b.buildActions, func(r *Reply) {
		r.ReplyType = TypeMessage
	})
	return b
}

//TypePhoto sets the type to TypePhoto
func (b *ReplyBuilder) TypePhoto() *ReplyBuilder {
	b.buildActions = append(b.buildActions, func(r *Reply) {
		r.ReplyType = TypePhoto
	})
	return b
}

//TypeAudio sets the type to TypeAudio
func (b *ReplyBuilder) TypeAudio() *ReplyBuilder {
	b.buildActions = append(b.buildActions, func(r *Reply) {
		r.ReplyType = TypeAudio
	})
	return b
}

//WithContent sets the content of the Reply message
func (b *ReplyBuilder) WithContent(content interface{}) *ReplyBuilder {
	b.buildActions = append(b.buildActions, func(r *Reply) {
		r.Content = content
	})
	return b
}

//ForChat sets the chatId
func (b *ReplyBuilder) ForChat(chatId string) *ReplyBuilder {
	b.buildActions = append(b.buildActions, func(r *Reply) {
		r.ChatId = chatId
	})
	return b
}

//Build builds the Reply object
func (b *ReplyBuilder) Build() Reply {
	emp := Reply{
		Content:   nil,
		ReplyType: -1,
		ChatId:    "0",
	}

	for _, a := range b.buildActions {
		a(&emp)
	}
	return emp
}
