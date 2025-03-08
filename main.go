package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

// Version set at compile-time
var Version string

func main() {
	// Load env-file if it exists first
	if filename, found := os.LookupEnv("PLUGIN_ENV_FILE"); found {
		_ = godotenv.Load(filename)
	}

	if _, err := os.Stat("/run/drone/env"); err == nil {
		_ = godotenv.Overload("/run/drone/env")
	}

	app := cli.NewApp()
	app.Name = "Drone Discord"
	app.Usage = "Sending message to Discord channel using Webhook"
	app.Copyright = "Copyright (c) " + strconv.Itoa(time.Now().Year()) + " Bo-Yi Wu"
	app.Authors = []*cli.Author{
		{
			Name:  "Bo-Yi Wu",
			Email: "appleboy.tw@gmail.com",
		},
	}
	app.Action = run
	app.Version = Version
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "webhook-url",
			Usage:   "discord webhook url",
			EnvVars: []string{"PLUGIN_WEBHOOK_URL", "WEBHOOK_URL", "DISCORD_WEBHOOK_URL", "INPUT_WEBHOOK_URL"},
		},
		&cli.StringFlag{
			Name:    "webhook-id",
			Usage:   "discord webhook id",
			EnvVars: []string{"PLUGIN_WEBHOOK_ID", "WEBHOOK_ID", "DISCORD_WEBHOOK_ID", "INPUT_WEBHOOK_ID"},
		},
		&cli.StringFlag{
			Name:    "webhook-token",
			Usage:   "discord webhook token",
			EnvVars: []string{"PLUGIN_WEBHOOK_TOKEN", "WEBHOOK_TOKEN", "DISCORD_WEBHOOK_TOKEN", "INPUT_WEBHOOK_TOKEN"},
		},
		&cli.StringSliceFlag{
			Name:    "message",
			Usage:   "the message contents (up to 2000 characters)",
			EnvVars: []string{"PLUGIN_MESSAGE", "DISCORD_MESSAGE", "MESSAGE", "INPUT_MESSAGE"},
		},
		&cli.StringSliceFlag{
			Name:    "file",
			Usage:   "the contents of the file being sent",
			EnvVars: []string{"PLUGIN_FILE", "DISCORD_FILE", "FILE", "INPUT_FILE"},
		},
		&cli.StringFlag{
			Name:    "color",
			Usage:   "color code of the embed",
			EnvVars: []string{"PLUGIN_COLOR", "COLOR", "INPUT_COLOR"},
		},
		&cli.BoolFlag{
			Name:    "wait",
			Usage:   "waits for server confirmation of message send before response, and returns the created message body",
			EnvVars: []string{"PLUGIN_WAIT", "WAIT", "INPUT_WAIT"},
		},
		&cli.BoolFlag{
			Name:    "tts",
			Usage:   "true if this is a TTS message",
			EnvVars: []string{"PLUGIN_TTS", "TTS", "INPUT_TTS"},
		},
		&cli.StringFlag{
			Name:    "username",
			Usage:   "override the default username of the webhook",
			EnvVars: []string{"PLUGIN_USERNAME", "USERNAME", "INPUT_USERNAME"},
		},
		&cli.StringFlag{
			Name:    "avatar-url",
			Usage:   "override the default avatar of the webhook",
			EnvVars: []string{"PLUGIN_AVATAR_URL", "AVATAR_URL", "INPUT_AVATAR_URL"},
		},
		&cli.BoolFlag{
			Name:    "drone",
			Usage:   "environment is drone",
			EnvVars: []string{"DRONE"},
		},
		&cli.StringFlag{
			Name:    "ci.environment",
			Usage:   "ci environment name",
			EnvVars: []string{"CI"},
		},
		&cli.StringFlag{
			Name:    "repo",
			Usage:   "repository owner and repository name",
			EnvVars: []string{"DRONE_REPO", "CI_REPO", "GITHUB_REPOSITORY"},
		},
		&cli.StringFlag{
			Name:    "repo.namespace",
			Usage:   "repository namespace",
			EnvVars: []string{"DRONE_REPO_OWNER", "DRONE_REPO_NAMESPACE", "CI_REPO_OWNER", "GITHUB_ACTOR"},
		},
		&cli.StringFlag{
			Name:    "repo.name",
			Usage:   "repository name",
			EnvVars: []string{"DRONE_REPO_NAME", "CI_REPO_NAME"},
		},
		&cli.StringFlag{
			Name:    "commit.sha",
			Usage:   "git commit sha",
			EnvVars: []string{"DRONE_COMMIT_SHA", "CI_COMMIT_SHA", "GITHUB_SHA"},
		},
		&cli.StringFlag{
			Name:    "commit.ref",
			Usage:   "git commit ref",
			EnvVars: []string{"DRONE_COMMIT_REF", "CI_COMMIT_REF", "GITHUB_REF"},
		},
		&cli.StringFlag{
			Name:    "commit.branch",
			Value:   "master",
			Usage:   "git commit branch",
			EnvVars: []string{"DRONE_COMMIT_BRANCH", "CI_COMMIT_BRANCH"},
		},
		&cli.StringFlag{
			Name:    "commit.link",
			Usage:   "git commit link",
			EnvVars: []string{"DRONE_COMMIT_LINK", "CI_PIPELINE_FORGE_URL"},
		},
		&cli.StringFlag{
			Name:    "commit.author",
			Usage:   "git author name",
			EnvVars: []string{"DRONE_COMMIT_AUTHOR", "CI_COMMIT_AUTHOR"},
		},
		&cli.StringFlag{
			Name:    "commit.author.email",
			Usage:   "git author email",
			EnvVars: []string{"DRONE_COMMIT_AUTHOR_EMAIL", "CI_COMMIT_AUTHOR_EMAIL"},
		},
		&cli.StringFlag{
			Name:    "commit.author.avatar",
			Usage:   "git author avatar",
			EnvVars: []string{"DRONE_COMMIT_AUTHOR_AVATAR", "CI_COMMIT_AUTHOR_AVATAR"},
		},
		&cli.StringFlag{
			Name:    "commit.message",
			Usage:   "commit message",
			EnvVars: []string{"DRONE_COMMIT_MESSAGE", "CI_COMMIT_MESSAGE"},
		},
		&cli.StringFlag{
			Name:    "source.branch",
			Value:   "develop",
			Usage:   "git source branch",
			EnvVars: []string{"DRONE_SOURCE_BRANCH", "CI_COMMIT_SOURCE_BRANCH"},
		},
		&cli.StringFlag{
			Name:    "build.event",
			Value:   "push",
			Usage:   "build event",
			EnvVars: []string{"DRONE_BUILD_EVENT", "CI_PIPELINE_EVENT"},
		},
		&cli.IntFlag{
			Name:    "build.number",
			Usage:   "build number",
			EnvVars: []string{"DRONE_BUILD_NUMBER", "CI_PIPELINE_NUMBER"},
		},
		&cli.StringFlag{
			Name:    "build.status",
			Usage:   "build status",
			Value:   "success",
			EnvVars: []string{"DRONE_BUILD_STATUS", "CI_PIPELINE_STATUS"},
		},
		&cli.StringFlag{
			Name:    "build.link",
			Usage:   "build link",
			EnvVars: []string{"DRONE_BUILD_LINK", "CI_PIPELINE_URL"},
		},
		&cli.StringFlag{
			Name:    "build.tag",
			Usage:   "build tag",
			EnvVars: []string{"DRONE_TAG", "CI_COMMIT_TAG"},
		},
		&cli.StringFlag{
			Name:    "pull.request",
			Usage:   "pull request",
			EnvVars: []string{"DRONE_PULL_REQUEST", "CI_COMMIT_PULL_REQUEST"},
		},
		&cli.Int64Flag{
			Name:    "build.started",
			Usage:   "build started",
			EnvVars: []string{"DRONE_BUILD_STARTED", "CI_PIPELINE_STARTED"},
		},
		&cli.Int64Flag{
			Name:    "build.finished",
			Usage:   "build finished",
			EnvVars: []string{"DRONE_BUILD_FINISHED", "CI_PIPELINE_FINISHED"},
		},
		&cli.BoolFlag{
			Name:    "github",
			Usage:   "Boolean value, indicates the runtime environment is GitHub Action.",
			EnvVars: []string{"PLUGIN_GITHUB", "GITHUB"},
		},
		&cli.StringFlag{
			Name:    "github.workflow",
			Usage:   "The name of the workflow.",
			EnvVars: []string{"GITHUB_WORKFLOW"},
		},
		&cli.StringFlag{
			Name:    "github.action",
			Usage:   "The name of the action.",
			EnvVars: []string{"GITHUB_ACTION"},
		},
		&cli.StringFlag{
			Name:    "github.event.name",
			Usage:   "The webhook name of the event that triggered the workflow.",
			EnvVars: []string{"GITHUB_EVENT_NAME"},
		},
		&cli.StringFlag{
			Name:    "github.event.path",
			Usage:   "The path to a file that contains the payload of the event that triggered the workflow. Value: /github/workflow/event.json.",
			EnvVars: []string{"GITHUB_EVENT_PATH"},
		},
		&cli.StringFlag{
			Name:    "github.workspace",
			Usage:   "The GitHub workspace path. Value: /github/workspace.",
			EnvVars: []string{"GITHUB_WORKSPACE"},
		},
		&cli.StringFlag{
			Name:    "deploy.to",
			Usage:   "Provides the target deployment environment for the running build. This value is only available to promotion and rollback pipelines.",
			EnvVars: []string{"DRONE_DEPLOY_TO", "CI_PIPELINE_DEPLOY_TARGET"},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		GitHub: GitHub{
			Workflow:  c.String("github.workflow"),
			Workspace: c.String("github.workspace"),
			Action:    c.String("github.action"),
			EventName: c.String("github.event.name"),
			EventPath: c.String("github.event.path"),
		},
		Repo: Repo{
			FullName:  c.String("repo"),
			Namespace: c.String("repo.namespace"),
			Name:      c.String("repo.name"),
		},
		Commit: Commit{
			Sha:     c.String("commit.sha"),
			Ref:     c.String("commit.ref"),
			Branch:  c.String("commit.branch"),
			Link:    c.String("commit.link"),
			Author:  c.String("commit.author"),
			Email:   c.String("commit.author.email"),
			Avatar:  c.String("commit.author.avatar"),
			Message: c.String("commit.message"),
		},
		Source: Source{
			Branch: c.String("source.branch"),
		},
		Build: Build{
			Tag:      c.String("build.tag"),
			Number:   c.Int("build.number"),
			Event:    c.String("build.event"),
			Status:   c.String("build.status"),
			Link:     c.String("build.link"),
			Started:  c.Int64("build.started"),
			Finished: c.Int64("build.finished"),
			PR:       c.String("pull.request"),
			DeployTo: c.String("deploy.to"),
		},
		Config: Config{
			webhookURL:   c.String("webhook-url"),
			WebhookID:    c.String("webhook-id"),
			WebhookToken: c.String("webhook-token"),
			Message:      c.StringSlice("message"),
			File:         c.StringSlice("file"),
			Color:        c.String("color"),
			Drone:        c.Bool("drone") || c.String("ci.environment") == "woodpecker",
			GitHub:       c.Bool("github"),
		},
		Payload: Payload{
			Wait:      c.Bool("wait"),
			Username:  c.String("username"),
			AvatarURL: c.String("avatar-url"),
			TTS:       c.Bool("tts"),
		},
	}

	return plugin.Exec(c.Context)
}
