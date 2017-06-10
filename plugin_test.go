package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingConfig(t *testing.T) {
	var plugin Plugin

	err := plugin.Exec()

	assert.NotNil(t, err)
}

func TestDefaultMessageFormat(t *testing.T) {
	plugin := Plugin{
		Repo: Repo{
			Name:  "go-hello",
			Owner: "appleboy",
		},
		Build: Build{
			Number:  101,
			Status:  "success",
			Link:    "https://github.com/appleboy/go-hello",
			Author:  "Bo-Yi Wu",
			Branch:  "master",
			Message: "update by drone line plugin.",
			Commit:  "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
		},
	}

	message := plugin.Message(plugin.Repo, plugin.Build)

	assert.Equal(t, []string{"[success] <https://github.com/appleboy/go-hello> (master)『update by drone line plugin.』by Bo-Yi Wu"}, message)
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
			Author:  "Bo-Yi Wu",
			Branch:  "master",
			Message: "update by drone discord plugin.",
			Commit:  "e7c4f0a63ceeb42a39ac7806f7b51f3f0d204fd2",
		},

		Config: Config{
			WebhookID:    os.Getenv("WEBHOOK_ID"),
			WebhookToken: os.Getenv("WEBHOOK_TOKEN"),
			Wait:         false,
			Message:      []string{"test one message from drone testing", "test two message from drone testing"},
			Username:     "drone-ci",
			TTS:          false,
		},
	}

	err := plugin.Exec()
	assert.Nil(t, err)

	plugin.Config.Message = []string{}
	err = plugin.Exec()
	assert.Nil(t, err)

	plugin.Config.Message = []string{"I am appleboy"}
	plugin.Config.TTS = true
	plugin.Config.Wait = true
	err = plugin.Exec()
	assert.Nil(t, err)
}
