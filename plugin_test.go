package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMissingConfig(t *testing.T) {
	plugin := Plugin{}

	err := plugin.Exec(context.Background())

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "missing discord config")
}

func TestSendPlainTextMessage(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			WebhookID:    os.Getenv("WEBHOOK_ID"),
			WebhookToken: os.Getenv("WEBHOOK_TOKEN"),
			Message:      []string{"Hello, world!", "This is a test."},
		},
		Payload: Payload{
			Username: "test-bot",
		},
	}

	err := plugin.Exec(context.Background())
	assert.NoError(t, err)
}

func TestSendEmbedMessage(t *testing.T) {
	plugin := Plugin{
		Config: Config{
			WebhookID:    os.Getenv("WEBHOOK_ID"),
			WebhookToken: os.Getenv("WEBHOOK_TOKEN"),
			Message:      []string{"This is an embed message."},
			Color:        "#48f442",
		},
		Payload: Payload{
			Username: "embed-bot",
		},
	}

	err := plugin.Exec(context.Background())
	assert.NoError(t, err)
}

func TestSendDefaultMessage(t *testing.T) {
	plugin := Plugin{
		Repo: Repo{
			Name:      "go-hello",
			Namespace: "appleboy",
		},
		Commit: Commit{
			Author:  "appleboy",
			Branch:  "master",
			Message: "feat: new feature",
			Avatar:  "https://avatars0.githubusercontent.com/u/21979?v=3&s=100",
		},
		Build: Build{
			Number: 101,
			Status: "success",
			Link:   "https://github.com/appleboy/go-hello",
			Event:  "push",
		},
		Config: Config{
			WebhookID:    os.Getenv("WEBHOOK_ID"),
			WebhookToken: os.Getenv("WEBHOOK_TOKEN"),
			Drone:        true,
		},
		Payload: Payload{
			Username: "default-bot",
		},
	}

	err := plugin.Exec(context.Background())
	assert.NoError(t, err)
}

func TestSendFile(t *testing.T) {
	// Create a dummy file for testing
	dummyFile, err := os.Create("test_file.txt")
	assert.NoError(t, err)
	_, err = dummyFile.WriteString("This is a test file.")
	assert.NoError(t, err)
	dummyFile.Close()
	defer os.Remove("test_file.txt")

	plugin := Plugin{
		Config: Config{
			WebhookID:    os.Getenv("WEBHOOK_ID"),
			WebhookToken: os.Getenv("WEBHOOK_TOKEN"),
			File:         []string{"test_file.txt"},
		},
		Payload: Payload{
			Username: "file-bot",
		},
	}

	err = plugin.Exec(context.Background())
	assert.NoError(t, err)
}

func TestColorConversion(t *testing.T) {
	tests := []struct {
		name         string
		colorHex     string
		expectedInt  int
		buildStatus  string
		expectedFall int
	}{
		{"valid hex", "#ffaa00", 16755200, "success", 1752220},
		{"invalid hex", "not-a-hex", 0, "failure", 16724530},
		{"status success", "", 0, "success", 1752220},
		{"status failure", "", 0, "failure", 16724530},
		{"status killed", "", 0, "killed", 16724530},
		{"status default", "", 0, "running", 16767280},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Plugin{
				Config: Config{Color: tt.colorHex},
				Build:  Build{Status: tt.buildStatus},
			}
			if tt.colorHex != "" {
				assert.Equal(t, tt.expectedInt, p.Color())
			} else {
				assert.Equal(t, tt.expectedFall, p.Color())
			}
		})
	}
}

func TestExecWithAllFeatures(t *testing.T) {
	time.Sleep(1 * time.Second)
	// Create a dummy file for testing
	dummyFile, err := os.Create("test_all.txt")
	assert.NoError(t, err)
	_, err = dummyFile.WriteString("This is a test file for a combined test.")
	assert.NoError(t, err)
	dummyFile.Close()
	defer os.Remove("test_all.txt")

	plugin := Plugin{
		Repo: Repo{
			Name:      "go-hello",
			Namespace: "appleboy",
		},
		Commit: Commit{
			Author:  "appleboy",
			Message: "Combined test with multiple features",
		},
		Build: Build{
			Status: "success",
			Link:   "http://example.com",
		},
		Config: Config{
			WebhookID:    os.Getenv("WEBHOOK_ID"),
			WebhookToken: os.Getenv("WEBHOOK_TOKEN"),
			Message:      []string{"First line of embed.", "Second line."},
			File:         []string{"test_all.txt"},
			Color:        "#32a852",
			Drone:        true,
		},
		Payload: Payload{
			Username: "super-bot",
		},
	}

	err = plugin.Exec(context.Background())
	assert.NoError(t, err)
}
