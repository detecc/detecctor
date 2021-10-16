package api

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type CommandBuilderTestSuite struct {
	suite.Suite
	builder *CommandBuilder
}

func (suite *CommandBuilderTestSuite) SetupTest() {
	suite.builder = NewCommandBuilder()
}

func (suite *CommandBuilderTestSuite) TestEmptyBuild() {
	suite.Equal(suite.builder.Build(), Command{
		Name:   "",
		Args:   []string{},
		ChatId: "0",
	})
}

func (suite *CommandBuilderTestSuite) TestLegitBuild() {
	cmd := suite.builder.WithName("name123").WithArgs([]string{"a", "b"}).FromChat("chatId123").Build()
	suite.Equal(cmd, Command{
		Name:   "/name123",
		Args:   []string{"a", "b"},
		ChatId: "chatId123",
	})
}

func (suite *CommandBuilderTestSuite) TestBuildWithoutArgs() {
	cmd := suite.builder.WithName("name123").FromChat("chatId123").Build()
	suite.Equal(cmd, Command{
		Name:   "/name123",
		Args:   []string{},
		ChatId: "chatId123",
	})
}

func TestCommandBuilder(t *testing.T) {
	suite.Run(t, new(CommandBuilderTestSuite))
}
