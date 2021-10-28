package api

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type MessageBuilderTestSuite struct {
	suite.Suite
	builder *MessageBuilder
}

func (suite *MessageBuilderTestSuite) SetupTest() {
	suite.builder = NewMessageBuilder()
}

func (suite *MessageBuilderTestSuite) TestEmptyBuild() {
	expected := ProxyMessage{
		ChatId:    "",
		Username:  "",
		Message:   "",
		MessageId: "",
	}
	suite.Equal(expected, suite.builder.Build())
}

func (suite *MessageBuilderTestSuite) TestLegitBuild() {
	expected := ProxyMessage{
		ChatId:    "chatId123",
		Username:  "user1",
		Message:   "Example message",
		MessageId: "Message123",
	}
	reply := suite.builder.WithId("chatId123").WithMessage("Message123", "Example message").FromUser("user1").Build()
	suite.Equal(expected, reply)
}

func (suite *MessageBuilderTestSuite) TestBuildWithoutOneAttribute() {
	expected := ProxyMessage{
		ChatId:    "",
		Username:  "",
		Message:   "Example message",
		MessageId: "Message123",
	}
	reply := suite.builder.WithMessage("Message123", "Example message").Build()
	suite.Equal(expected, reply)
}

func TestMessageBuilder(t *testing.T) {
	suite.Run(t, new(MessageBuilderTestSuite))
}
