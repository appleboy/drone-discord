package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/appleboy/drone-template-lib/template"
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

	// Repo information.
	Repo struct {
		FullName  string
		Namespace string
		Name      string
	}

	// Commit information.
	Commit struct {
		Sha     string
		Ref     string
		Branch  string
		Link    string
		Author  string
		Avatar  string
		Email   string
		Message string
	}

	Source struct {
		Branch string
	}

	// Build information.
	Build struct {
		Tag      string
		Event    string
		Number   int
		Status   string
		Link     string
		Started  int64
		Finished int64
		PR       string
		DeployTo string
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
		Username  string        `json:"username,omitempty"`
		AvatarURL string        `json:"avatar_url"`
		TTS       bool          `json:"tts"`
		Embeds    []EmbedObject `json:"embeds"`
	}

	// Plugin values.
	Plugin struct {
		GitHub  GitHub
		Repo    Repo
		Build   Build
		Source  Source
		Config  Config
		Payload Payload
		Commit  Commit
	}
)

func templateMessage(t string, plugin Plugin) (string, error) {
	return template.RenderTrim(t, plugin)
}

// Creates a new file upload http request with optional extra params
// https://matt.aimonetti.net/posts/2013/07/01/golang-multipart-file-upload-example/
func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, err
}

// Exec executes the plugin.
func (p *Plugin) Exec(ctx context.Context) error {
	if p.Config.WebhookID == "" || p.Config.WebhookToken == "" {
		return errors.New("missing discord config")
	}

	if len(p.Config.Message) == 0 {
		object := p.Template()
		p.Payload.Embeds = []EmbedObject{object}
		err := p.SendMessage(ctx)
		if err != nil {
			return err
		}
	}

	if len(p.Config.Message) > 0 {
		for _, m := range p.Config.Message {
			txt, err := templateMessage(m, *p)
			if err != nil {
				return err
			}

			if len(p.Config.Color) != 0 {
				object := p.DefaultTemplate(txt)
				p.Payload.Embeds = append(p.Payload.Embeds, object)
			} else {
				p.Payload.Content = txt
				err = p.SendMessage(ctx)
				if err != nil {
					return err
				}
			}
		}

		if len(p.Payload.Embeds) > 0 {
			err := p.SendMessage(ctx)
			if err != nil {
				return err
			}
		}
	}

	for _, f := range p.Config.File {
		if f == "" {
			continue
		}
		err := p.SendFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

// SendFile upload file to discord
func (p *Plugin) SendFile(file string) error {
	webhookURL := fmt.Sprintf("https://discord.com/api/webhooks/%s/%s", p.Config.WebhookID, p.Config.WebhookToken)
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
func (p *Plugin) SendMessage(ctx context.Context) error {
	webhookURL := fmt.Sprintf("https://discordapp.com/api/webhooks/%s/%s", p.Config.WebhookID, p.Config.WebhookToken)
	b := new(bytes.Buffer)
	if err := json.NewEncoder(b).Encode(p.Payload); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhookURL, b)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
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
	var description string
	switch {
	case p.Config.GitHub:
		description = fmt.Sprintf("%s/%s triggered by %s (%s)",
			p.Repo.FullName,
			p.GitHub.Workflow,
			p.Repo.Namespace,
			p.GitHub.EventName,
		)
	case p.Build.Event == "push":
		description = fmt.Sprintf("%s pushed to %s", p.Commit.Author, p.Commit.Branch)
	case p.Build.Event == "pull_request":
		branch := p.Commit.Ref
		if branch == "" {
			branch = p.Commit.Branch
		}
		description = fmt.Sprintf("%s updated pull request %s", p.Commit.Author, branch)
	case p.Build.Event == "tag":
		description = fmt.Sprintf("%s pushed tag %s", p.Commit.Author, p.Commit.Branch)
	}

	return EmbedObject{
		Title:       p.Commit.Message,
		Description: description,
		URL:         p.Build.Link,
		Color:       p.Color(),
		Author: EmbedAuthorObject{
			Name:    p.Commit.Author,
			IconURL: p.Commit.Avatar,
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
