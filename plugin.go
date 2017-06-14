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
		Message      []string
	}

	// EmbedFooterObject for Embed Footer Structure.
	EmbedFooterObject struct {
		Text string `json:"text"`
	}

	// EmbedAuthorObject for Embed Author Structure
	EmbedAuthorObject struct {
		Name    string `json:"name"`
		URL     string `json:"url"`
		IconURL string `json:"icon_url"`
	}

	// EmbedFieldObject for Embed Field Structure
	EmbedFieldObject struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	}

	// EmbedObject is for Embed Structure
	EmbedObject struct {
		Title       string              `json:"title"`
		Description string              `json:"description"`
		URL         string              `json:"url"`
		Color       int                 `json:"color"`
		Footer      *EmbedFooterObject  `json:"footer"`
		Author      *EmbedAuthorObject  `json:"author"`
		Fields      []*EmbedFieldObject `json:"fields"`
	}

	// Payload struct
	Payload struct {
		Wait      bool           `json:"wait"`
		Content   string         `json:"content"`
		Username  string         `json:"username"`
		AvatarURL string         `json:"avatar_url"`
		TTS       bool           `json:"tts"`
		Embeds    []*EmbedObject `json:"embeds"`
	}

	// Plugin values.
	Plugin struct {
		Repo    Repo
		Build   Build
		Config  Config
		Payload Payload
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
	if len(p.Config.Message) > 0 {
		messages = p.Config.Message
	} else {
		messages = p.Message(p.Repo, p.Build)
	}

	for _, m := range messages {
		txt, err := template.RenderTrim(m, p)
		if err != nil {
			return err
		}

		// update content
		p.Payload.Content = txt
		b := new(bytes.Buffer)
		json.NewEncoder(b).Encode(p.Payload)
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
