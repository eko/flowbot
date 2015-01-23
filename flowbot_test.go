// Flowbot - A Flowdock robot library written in Go
//
// Author: Vincent Composieux <vincent.composieux@gmail.com>

package flowbot

import (
    "testing"
)

// Tests pre-defined and variables values
func TestVariablesValues(t *testing.T) {
    if (FlowdockRobotName != "Flowbot") {
        t.Error("Should have the default robot name: Flowbot")
    }

    FlowdockRobotName    = "TestBotName"

    if (FlowdockRobotName != "TestBotName") {
        t.Error("Should have the default robot name: Flowbot")
    }

    FlowdockChatUrl      = "https://api.flowdock.com/v1/messages/chat/"
    FlowdockInboxUrl     = "https://api.flowdock.com/v1/messages/team_inbox/"

    FlowdockStreamUrl    = "http://test.url.com"
    FlowdockFlowToken    = "test-token"
    FlowdockAuthUsername = "toto"
    FlowdockAuthPassword = "password"

    if (FlowdockStreamUrl != "http://test.url.com" || FlowdockFlowToken != "test-token" || FlowdockAuthUsername != "toto" || FlowdockAuthPassword != "password") {
        t.Error("Should have the URL configured")
    }
}

// Tests the command addition
func TestAddCommand(t *testing.T) {
    if (len(commands) != 0) {
        t.Error("Should have 0 commands by default")
    }

    AddCommand("test-a", func (command Command, entry Entry) {})
    AddCommand("test-b", func (command Command, entry Entry) {})

    if (len(commands) != 2) {
        t.Error("Should have 2 commands added")
    }
}
