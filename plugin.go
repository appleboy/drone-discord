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
		Avatar   string
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
		Title       string             `json:"title"`
		Description string             `json:"description"`
		URL         string             `json:"url"`
		Color       int                `json:"color"`
		Footer      EmbedFooterObject  `json:"footer"`
		Author      EmbedAuthorObject  `json:"author"`
		Fields      []EmbedFieldObject `json:"fields"`
	}

	// Payload struct
	Payload struct {
		Wait      bool          `json:"wait"`
		Content   string        `json:"content"`
		Username  string        `json:"username"`
		AvatarURL string        `json:"avatar_url"`
		TTS       bool          `json:"tts"`
		Embeds    []EmbedObject `json:"embeds"`
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
func (p *Plugin) Exec() error {
	if len(p.Config.WebhookID) == 0 || len(p.Config.WebhookToken) == 0 {
		log.Println("missing discord config")

		return errors.New("missing discord config")
	}

	if len(p.Config.Message) > 0 {
		for _, m := range p.Config.Message {
			txt, err := template.RenderTrim(m, p)
			if err != nil {
				return err
			}

			// update content
			p.Payload.Content = txt
			err = p.Send()
			if err != nil {
				return err
			}
		}
		return nil
	}

	// set default message
	p.Message()
	err := p.Send()
	if err != nil {
		return err
	}

	return nil
}

// Send discord message.
func (p *Plugin) Send() error {
	webhookURL := fmt.Sprintf("https://discordapp.com/api/webhooks/%s/%s", p.Config.WebhookID, p.Config.WebhookToken)
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(p.Payload)
	_, err := http.Post(webhookURL, "application/json; charset=utf-8", b)

	if err != nil {
		return err
	}

	return nil
}

// Message is plugin default message.
func (p *Plugin) Message() {
	p.Payload.Embeds = []EmbedObject{
		{
			Title: p.Build.Message,
			URL:   p.Build.Link,
			Color: color(p.Build),
			Author: EmbedAuthorObject{
				Name:    p.Build.Author,
				IconURL: p.Build.Avatar,
			},
			Footer: EmbedFooterObject{
				Text: "Power by Drone Discord Plugin",
			},
		},
	}
}

func color(build Build) int {
	switch build.Status {
	case "success":
		// #1ac600
		return 1754624
	case "failure", "error", "killed":
		// #ff3232
		return 16724530
	default:
		// #ffd930
		return 16767280
	}
}
