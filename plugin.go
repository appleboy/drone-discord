package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/appleboy/drone-facebook/template"
)

const (
	// DroneIconURL default drone logo url
	DroneIconURL = "https://c1.staticflickr.com/5/4236/34957940160_435d83114f_z.jpg"
	// DroneDesc default drone description
	DroneDesc = "Powered by Drone Discord Plugin"
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
		RefSpec  string
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
		Color        string
		Message      []string
		Drone        bool
	}

	// EmbedFooterObject for Embed Footer Structure.
	EmbedFooterObject struct {
		Text    string `json:"text"`
		IconURL string `json:"icon_url"`
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

	if p.Config.Drone && len(p.Config.Message) == 0 {
		object := p.DroneTemplate()
		p.Payload.Embeds = []EmbedObject{object}
		err := p.Send()
		if err != nil {
			return err
		}
	}

	if len(p.Config.Message) > 0 {
		for _, m := range p.Config.Message {
			txt, err := template.RenderTrim(m, p)
			if err != nil {
				return err
			}

			if len(p.Config.Color) != 0 {
				object := p.DefaultTemplate(txt)
				p.Payload.Embeds = append(p.Payload.Embeds, object)
			} else {
				p.Payload.Content = txt
				err = p.Send()
				if err != nil {
					return err
				}
			}
		}

		if len(p.Payload.Embeds) > 0 {
			err := p.Send()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Send discord message.
func (p *Plugin) Send() error {
	webhookURL := fmt.Sprintf("https://discordapp.com/api/webhooks/%s/%s", p.Config.WebhookID, p.Config.WebhookToken)
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(p.Payload); err != nil {
		return err
	}
	_, err := http.Post(webhookURL, "application/json; charset=utf-8", b)

	if err != nil {
		return err
	}

	return nil
}

// DefaultTemplate is plugin default template for Drone CI.
func (p *Plugin) DefaultTemplate(title string) EmbedObject {
	return EmbedObject{
		Title: title,
		Color: p.Color(),
	}
}

// DroneTemplate is plugin default template for Drone CI.
func (p *Plugin) DroneTemplate() EmbedObject {
	description := ""
	switch p.Build.Event {
	case "push":
		description = fmt.Sprintf("%s pushed to %s", p.Build.Author, p.Build.Branch)
	case "pull_request":
		branch := ""
		if p.Build.RefSpec != "" {
			branch = p.Build.RefSpec
		} else {
			branch = p.Build.Branch
		}
		description = fmt.Sprintf("%s updated pull request %s", p.Build.Author, branch)
	case "tag":
		description = fmt.Sprintf("%s pushed tag %s", p.Build.Author, p.Build.Branch)
	}

	return EmbedObject{
		Title:       p.Build.Message,
		Description: description,
		URL:         p.Build.Link,
		Color:       p.Color(),
		Author: EmbedAuthorObject{
			Name:    p.Build.Author,
			IconURL: p.Build.Avatar,
		},
		Footer: EmbedFooterObject{
			Text:    DroneDesc,
			IconURL: DroneIconURL,
		},
	}
}

// Color code of the embed
func (p *Plugin) Color() int {
	if p.Config.Color != "" {
		p.Config.Color = strings.Replace(p.Config.Color, "#", "", -1)
		if s, err := strconv.ParseInt(p.Config.Color, 16, 32); err == nil {
			return int(s)
		}
	}

	switch p.Build.Status {
	case "success":
		// #1ac600 green
		return 1754624
	case "failure", "error", "killed":
		// #ff3232 red
		return 16724530
	default:
		// #ffd930 yellow
		return 16767280
	}
}
