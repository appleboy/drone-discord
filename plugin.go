package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/appleboy/drone-facebook/template"
)

type (
	// Repo information
	Repo struct {
		Owner string
		Name  string
	}

	// Build information
	Build struct {
		Tag      string
		Event    string
		Number   int
		Commit   string
		Branch   string
		Author   string
		Message  string
		Email    string
		Status   string
		Link     string
		Started  float64
		Finished float64
	}

	// Config for the plugin.
	Config struct {
		WebhookID    string
		WebhookToken string
		Wait         bool
		Content      []string
		Username     string
		AvatarURL    string
		TTS          bool
	}

	// Plugin values.
	Plugin struct {
		Repo   Repo
		Build  Build
		Config Config
	}
)

// Exec executes the plugin.
func (p Plugin) Exec() error {
	if len(p.Config.WebhookID) == 0 || len(p.Config.WebhookToken) == 0 {
		log.Println("missing discord config")

		return errors.New("missing discord config")
	}

	webhookURL := fmt.Sprintf("https://discordapp.com/api/webhooks/%s/%s", p.Config.WebhookID, p.Config.WebhookToken)

	var messages []string
	if len(p.Config.Content) > 0 {
		messages = p.Config.Content
	} else {
		messages = p.Message(p.Repo, p.Build)
	}

	for _, m := range messages {
		txt, err := template.RenderTrim(m, p)
		if err != nil {
			return err
		}

		content := map[string]interface{}{
			"wait":       p.Config.Wait,
			"content":    txt,
			"username":   p.Config.Username,
			"avatar_url": p.Config.AvatarURL,
			"tts":        p.Config.TTS,
		}
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(content)
		_, err = http.Post(webhookURL, "application/json; charset=utf-8", b)

		if err != nil {
			return err
		}
	}

	return nil
}

// Message is plugin default message.
func (p Plugin) Message(repo Repo, build Build) []string {
	return []string{fmt.Sprintf("[%s] <%s> (%s)『%s』by %s",
		build.Status,
		build.Link,
		build.Branch,
		build.Message,
		build.Author,
	)}
}
