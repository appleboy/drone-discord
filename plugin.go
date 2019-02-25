package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
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
	// GitHub information.
	GitHub struct {
		Workflow  string
		Workspace string
		Action    string
		EventName string
		EventPath string
	}

	// Repo information
	Repo struct {
		FullName  string
		Namespace string
		Name      string
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
		File         []string
		Drone        bool
		GitHub       bool
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
		GitHub  GitHub
		Repo    Repo
		Build   Build
		Config  Config
		Payload Payload
	}
)

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	fileContents, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}
	file.Close()

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, fi.Name())
	if err != nil {
		return nil, err
	}
	if _, err := part.Write(fileContents); err != nil {
		return nil, err
	}

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest("POST", uri, body)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	return request, err
}

// Exec executes the plugin.
func (p *Plugin) Exec() error {
	if p.Config.WebhookID == "" || p.Config.WebhookToken == "" {
		return errors.New("missing discord config")
	}

	if len(p.Config.Message) == 0 {
		object := p.Template()
		p.Payload.Embeds = []EmbedObject{object}
		err := p.SendMessage()
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
				err = p.SendMessage()
				if err != nil {
					return err
				}
			}
		}

		if len(p.Payload.Embeds) > 0 {
			err := p.SendMessage()
			if err != nil {
				return err
			}
		}
	}

	for _, f := range p.Config.File {
		err := p.SendFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

// SendFile upload file to discord
func (p *Plugin) SendFile(file string) error {
	webhookURL := fmt.Sprintf("https://discordapp.com/api/webhooks/%s/%s", p.Config.WebhookID, p.Config.WebhookToken)
	extraParams := map[string]string{}

	if p.Payload.Username != "" {
		extraParams["username"] = p.Payload.Username
	}

	if p.Payload.AvatarURL != "" {
		extraParams["avatar_url"] = p.Payload.AvatarURL
	}

	if p.Payload.TTS {
		extraParams["tts"] = "true"
	}

	request, err := newfileUploadRequest(
		webhookURL,
		extraParams,
		"file",
		file,
	)
	if err != nil {
		return err
	}
	client := &http.Client{}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}

// SendMessage to send discord message.
func (p *Plugin) SendMessage() error {
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

// Template is plugin default template for Drone CI or GitHub Action.
func (p *Plugin) Template() EmbedObject {
	if p.Config.GitHub {
		message := fmt.Sprintf("%s/%s triggered by %s (%s)",
			p.Repo.FullName,
			p.GitHub.Workflow,
			p.Repo.Namespace,
			p.GitHub.EventName,
		)

		return EmbedObject{
			Title: message,
			Color: p.Color(),
			Footer: EmbedFooterObject{
				Text:    DroneDesc,
				IconURL: DroneIconURL,
			},
		}
	}

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

// Clear reset to default
func (p *Plugin) Clear() {
	// clear content field.
	p.Payload.Content = ""
	p.Payload.Embeds = []EmbedObject{}
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
		// green
		return 0x1ac600
	case "failure", "error", "killed":
		// red
		return 0xff3232
	default:
		// yellow
		return 0xffd930
	}
}
