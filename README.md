# drone-[discord](https://discordapp.com)

![logo](images/discord-logo.png)

Drone plugin for sending message to Discord channel using Webhook.

[![GoDoc](https://godoc.org/github.com/appleboy/drone-discord?status.svg)](https://godoc.org/github.com/appleboy/drone-discord)
[![Build Status](https://cloud.drone.io/api/badges/appleboy/drone-discord/status.svg)](https://cloud.drone.io/appleboy/drone-discord)
[![Build status](https://ci.appveyor.com/api/projects/status/xj24ye9lu68a9sqm?svg=true)](https://ci.appveyor.com/project/appleboy/drone-discord-bne7m)
[![codecov](https://codecov.io/gh/appleboy/drone-discord/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-discord)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-discord)](https://goreportcard.com/report/github.com/appleboy/drone-discord)
[![Docker Pulls](https://img.shields.io/docker/pulls/appleboy/drone-discord.svg)](https://hub.docker.com/r/appleboy/drone-discord/)
[![microbadger](https://images.microbadger.com/badges/image/appleboy/drone-discord:linux-amd64.svg)](https://microbadger.com/images/appleboy/drone-discord:linux-amd64 "Get your own image badge on microbadger.com")

Webhooks are a low-effort way to post messages to channels in Discord. They do not require a bot user or authentication to use. See more [api document information](https://discordapp.com/developers/docs/resources/webhook). For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/appleboy/drone-discord/).

Sending discord message using a binary, docker or [Drone CI](http://docs.drone.io/).

## Features

* [x] Send Multiple Messages
* [x] Send Multiple Files

## Build or Download a binary

The pre-compiled binaries can be downloaded from [release page](https://github.com/appleboy/drone-discord/releases). Support the following OS type.

* Windows amd64/386
* Linux arm/amd64/386
* Darwin amd64/386

With `Go` installed

```sh
go get -u -v github.com/appleboy/drone-discord
```

or build the binary with the following command:

```sh
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go test -cover ./...

go build -v -a -tags netgo -o release/linux/amd64/drone-discord
```

## Usage

There are three ways to send notification.

### Usage from binary

```bash
drone-discord \
  --webhook-id xxxx \
  --webhook-token xxxx \
  --message "Test Message"
```

### Usage from docker

```bash
docker run --rm \
  -e WEBHOOK_ID=xxxxxxx \
  -e WEBHOOK_TOKEN=xxxxxxx \
  -e WAIT=false \
  -e TTS=false \
  -e USERNAME=test \
  -e AVATAR_URL=http://example.com/xxxx.png \
  -e MESSAGE=test \
  appleboy/drone-discord
```

### Usage from drone ci

#### Send Notification

Execute from the working directory:

```sh
docker run --rm \
  -e WEBHOOK_ID=xxxxxxx \
  -e WEBHOOK_TOKEN=xxxxxxx \
  -e WAIT=false \
  -e TTS=false \
  -e USERNAME=test \
  -e AVATAR_URL=http://example.com/xxxx.png \
  -e MESSAGE=test \
  -e DRONE_REPO_OWNER=appleboy \
  -e DRONE_REPO_NAME=go-hello \
  -e DRONE_COMMIT_SHA=e5e82b5eb3737205c25955dcc3dcacc839b7be52 \
  -e DRONE_COMMIT_BRANCH=master \
  -e DRONE_COMMIT_AUTHOR=appleboy \
  -e DRONE_COMMIT_AUTHOR_EMAIL=appleboy@gmail.com \
  -e DRONE_COMMIT_MESSAGE=Test_Your_Commit \
  -e DRONE_BUILD_NUMBER=1 \
  -e DRONE_BUILD_STATUS=success \
  -e DRONE_BUILD_LINK=http://github.com/appleboy/go-hello \
  -e DRONE_JOB_STARTED=1477550550 \
  -e DRONE_JOB_FINISHED=1477550750 \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  appleboy/drone-discord
```

## Declarative environment config usage

You can get more [information](DOCS.md) about how to use this plugin in drone.

```yml
- name: msg status
  image: appleboy/drone-discord
  settings:
    webhook_id:
      from_secret: discord_id
    webhook_token:
      from_secret: discord_token
    message: "{{#success build.status}}✅{{else}}❌{{/success}}  Repository `[{{repo.name}}/{{commit.branch}}]` triggered by event `[{{uppercase build.event}}]` for build.\n    - Commit [[{{commit.sha}}]({{commit.link}})]\n    - Author `[{{commit.author}} / {{commit.email}}]`\n    - Message: {{commit.message}}    - Drone build [[#{{build.number}}]({{build.link}})] reported `[{{uppercase build.status}}]` at `[{{datetime build.finished \"2006.01.02 15:04\" \"\"}}]`\n"
    when:
      status: [ success, failure, changed ]
```

```yml
- name: multi line msg status 
  ...
    message: >
      Line one
      Line two
```

## Testing

Test the package with the following command:

```sh
make test
```
