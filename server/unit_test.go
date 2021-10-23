package server

import (
	"github.com/detecc/detecctor/cache"
	"github.com/detecc/detecctor/shared"
	cache2 "github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ServerFuncTestSuite struct {
	suite.Suite
	cache *cache2.Cache
}

func (suite *ServerFuncTestSuite) SetupTest() {
	cache.Cache = cache2.New(time.Minute, time.Minute)
	suite.cache = cache.Memory()
}

func (suite *ServerFuncTestSuite) TestNodesAndCommands() {
	GenerateChatAuthenticationToken("chatId123")
	_, isFound := suite.cache.Get("auth-token-chatId123")
	suite.True(isFound)
}

func (suite *ServerFuncTestSuite) TestGeneratePayloadId() {
	payload := shared.Payload{
		Id:             "",
		ServiceNodeKey: "snKey123",
		Data:           "data",
		Command:        "/cmd",
		Success:        true,
		Error:          "",
	}

	generatePayloadId(&payload, "chatId1")
	suite.NotEqual("", payload.Id)

	item, isFound := suite.cache.Get(payload.Id)

	suite.True(isFound)
	suite.EqualValues("chatId1", item.(string))
}

func (suite *ServerFuncTestSuite) TestInterpretSubscriptionCommand() {
	expected := map[string]string{"nodes": "node1,node2,node3", "commands": "/auth,/get_status"}
	actual := interpretSubscriptionCommand([]string{"nodes=node1,node2,node3", "commands=/auth,/get_status"})
	suite.EqualValues(expected, actual)
}

func (suite *ServerFuncTestSuite) TestGetNodesAndCommands() {
	expectedNodes := []string{"node1", "node2", "node3"}
	expectedCommands := []string{"/auth", "/get_status"}
	inputMap := map[string]string{"nodes": "node1,node2,node3", "commands": "/auth,/get_status"}

	nodes, commands := getNodesAndCommands(inputMap)

	suite.EqualValues(expectedNodes, nodes)
	suite.EqualValues(expectedCommands, commands)
}

func TestServerSuite(t *testing.T) {
	suite.Run(t, new(ServerFuncTestSuite))
}
