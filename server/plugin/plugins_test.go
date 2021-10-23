package plugin

import (
	"github.com/detecc/detecctor/bot/api"
	"github.com/detecc/detecctor/shared"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type PluginMock struct {
	mock.Mock
	Handler
}

func (p *PluginMock) Response(payload shared.Payload) api.Reply {
	args := p.Called(payload)
	return args.Get(0).(api.Reply)
}

func (p *PluginMock) Execute(args ...string) ([]shared.Payload, error) {
	args2 := p.Called(args)
	return args2.Get(0).([]shared.Payload), args2.Error(0)
}

func (p *PluginMock) GetMetadata() Metadata {
	args := p.Called()
	return args.Get(0).(Metadata)
}

type PluginManagerTestSuite struct {
	suite.Suite
	pluginManager *PluginManager
	pluginMock    PluginMock
}

func (suite *PluginManagerTestSuite) SetupTest() {
	suite.pluginManager = GetPluginManager()
	suite.pluginMock = PluginMock{}
	suite.pluginManager.AddPlugin("/example", &suite.pluginMock)
}

func (suite *PluginManagerTestSuite) TestGetPlugin() {
	plugin, err := suite.pluginManager.GetPlugin("/example")
	suite.NoError(err)
	suite.EqualValues(&suite.pluginMock, plugin)

	plugin, err = suite.pluginManager.GetPlugin("/example1")
	suite.Error(err)
	suite.Nil(plugin)
}

func (suite *PluginManagerTestSuite) TestAddPlugin() {
	suite.pluginManager.AddPlugin("/example", &suite.pluginMock)

	plugin, err := suite.pluginManager.GetPlugin("/example")
	suite.NoError(err)
	suite.EqualValues(&suite.pluginMock, plugin)
}

func (suite *PluginManagerTestSuite) TestHasPlugin() {
	suite.True(suite.pluginManager.HasPlugin("/example"))

	suite.False(suite.pluginManager.HasPlugin("/example1"))
}

func TestGetPluginManager(t *testing.T) {
	suite.Run(t, new(PluginManagerTestSuite))
}
