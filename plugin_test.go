package main

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMissingConfig(t *testing.T) {
	var plugin Plugin

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestSendMessage(t *testing.T) {
	plugin := Plugin{
		Repo: Repo{
			Name:  "go-hello",
			Owner: "appleboy",
		},
		Build: Build{
			Number:  101,
			Status:  "success",
			Link:    "https://github.com/appleboy/go-hello",
			Author:  "appleboy",
			Branch:  "master",
			Event:   "tag",
			Message: "update by drone discord plugin. \r\n update by drone discord plugin.",
			Commit:  "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
			Avatar:  "https://avatars0.githubusercontent.com/u/21979?v=3&s=100",
		},

		Config: Config{
			WebhookID:    os.Getenv("WEBHOOK_ID"),
			WebhookToken: os.Getenv("WEBHOOK_TOKEN"),
			Message:      []string{"test one message from drone testing", "test two message from drone testing"},
			Drone:        true,
		},

		Payload: Payload{
			Username: "drone-ci",
			TTS:      false,
			Wait:     false,
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)

	plugin.Config.Message = []string{"I am appleboy"}
	plugin.Payload.TTS = true
	plugin.Payload.Wait = true
	err = plugin.Exec()
	assert.Nil(t, err)

	// send success embed message
	plugin.Config.Message = []string{}
	err = plugin.Exec()
	assert.Nil(t, err)

	// send success embed message
	plugin.Build.Status = "failure"
	plugin.Build.Message = "send failure embed message"
	err = plugin.Exec()
	assert.Nil(t, err)
	time.Sleep(1 * time.Second)

	// send default embed message
	plugin.Build.Status = "test"
	plugin.Build.Message = "send default embed message"
	err = plugin.Exec()
	assert.Nil(t, err)

	//change color for embed message
	plugin.Config.Color = "#4842f4"
	plugin.Build.Message = "Change embed color to #4842f4"
	err = plugin.Exec()
	assert.Nil(t, err)
}
