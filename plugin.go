package main

import (
	"fmt"
	"net/http"
	"strings"
)

type (
	// Repo information
	Repo struct {
		Owner string
		Name  string
	}

	// Build information
	Build struct {
		Tag    string
		Event  string
		Number int
		Commit string
		Branch string
		Author string
		Status string
		Link   string
	}

	// Config for the plugin.
	Config struct {
		WebhookID    string
		WebhookToken string
		Wait         bool
		Content      string
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
	attachment := slack.Attachment{
		Text:       message(p.Repo, p.Build),
		Fallback:   fallback(p.Repo, p.Build),
		Color:      color(p.Build),
		MarkdownIn: []string{"text", "fallback"},
		ImageURL:   p.Config.ImageURL,
	}

	payload := slack.WebHookPostPayload{}
	payload.Username = p.Config.Username
	payload.Attachments = []*slack.Attachment{&attachment}
	payload.IconUrl = p.Config.IconURL
	payload.IconEmoji = p.Config.IconEmoji

	if p.Config.Recipient != "" {
		payload.Channel = prepend("@", p.Config.Recipient)
	} else if p.Config.Channel != "" {
		payload.Channel = prepend("#", p.Config.Channel)
	}

	if p.Config.Template != "" {
		txt, err := RenderTrim(p.Config.Template, p)
		if err != nil {
			return err
		}
		attachment.Text = txt
	}

	client := slack.NewWebHook(p.Config.Webhook)
	return client.PostMessage(&payload)
}

func message(repo Repo, build Build) string {
	return fmt.Sprintf("*%s* <%s|%s/%s#%s> (%s) by %s",
		build.Status,
		build.Link,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func fallback(repo Repo, build Build) string {
	return fmt.Sprintf("%s %s/%s#%s (%s) by %s",
		build.Status,
		repo.Owner,
		repo.Name,
		build.Commit[:8],
		build.Branch,
		build.Author,
	)
}

func prepend(prefix, s string) string {
	if !strings.HasPrefix(s, prefix) {
		return prefix + s
	}
	return s
}
