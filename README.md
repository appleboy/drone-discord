<img src="line.png">

# drone-[discord](https://discordapp.com)

Drone plugin for sending message to Discord channel using Webhook.

[![GoDoc](https://godoc.org/github.com/appleboy/drone-discord?status.svg)](https://godoc.org/github.com/appleboy/drone-discord)
[![Build Status](http://drone.wu-boy.com/api/badges/appleboy/drone-discord/status.svg)](http://drone.wu-boy.com/appleboy/drone-discord)
[![codecov](https://codecov.io/gh/appleboy/drone-discord/branch/master/graph/badge.svg)](https://codecov.io/gh/appleboy/drone-discord)
[![Go Report Card](https://goreportcard.com/badge/github.com/appleboy/drone-discord)](https://goreportcard.com/report/github.com/appleboy/drone-discord)
[![Docker Pulls](https://img.shields.io/docker/pulls/appleboy/drone-discord.svg)](https://hub.docker.com/r/appleboy/drone-discord/)
[![](https://images.microbadger.com/badges/image/appleboy/drone-discord.svg)](https://microbadger.com/images/appleboy/drone-discord "Get your own image badge on microbadger.com")
[![Release](https://github-release-version.herokuapp.com/github/appleboy/drone-discord/release.svg?style=flat)](https://github.com/appleboy/drone-discord/releases/latest)

Webhooks are a low-effort way to post messages to channels in Discord. They do not require a bot user or authentication to use. See more [api document information](https://discordapp.com/developers/docs/resources/webhook).

Sending discord message using a binary, docker or [Drone CI](http://docs.drone.io/).

## Build or Download a binary

The pre-compiled binaries can be downloaded from [release page](https://github.com/appleboy/drone-discord/releases). Support the following OS type.

* Windows amd64/386
* Linux amd64/386
* Darwin amd64/386

With `Go` installed

```
$ go get -u -v github.com/appleboy/drone-discord
``` 

or build the binary with the following command:

```
$ make build
```

## Docker

Build the docker image with the following commands:

```
$ make docker
```

Please note incorrectly building the image for the correct x64 linux and with
CGO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-discord' not found or does not exist..
```

## Usage

There are three ways to send notification.

* [usage from binary](#usage-from-binary)
* [usage from docker](#usage-from-docker)
* [usage from drone ci](#usage-from-drone-ci)

<a name="usage-from-binary"></a>
### Usage from binary

#### Send Notification

```bash
drone-discord \
  --webhook-id xxxx \
  --webhook-token xxxx \
  --content "Test Message"
```

<a name="usage-from-docker"></a>
### Usage from docker

#### Send Notification

```bash
docker run --rm \
  -e WEBHOOK_ID=xxxxxxx \
  -e WEBHOOK_TOKEN=xxxxxxx \
  -e WAIT=false \
  -e TTS=false \
  -e USERNAME=test \
  -e AVATAR_URL=http://example.com/xxxx.png \
  -e CONTENT=test \
  appleboy/drone-discord
```

<a name="usage-from-drone-ci"></a>
### Usage from drone ci

#### Send Notification

Execute from the working directory:

```bash
docker run --rm \
  -e WEBHOOK_ID=xxxxxxx \
  -e WEBHOOK_TOKEN=xxxxxxx \
  -e WAIT=false \
  -e TTS=false \
  -e USERNAME=test \
  -e AVATAR_URL=http://example.com/xxxx.png \
  -e CONTENT=test \
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

You can get more [information](DOCS.md) about how to use scp plugin in drone.

## Testing

Test the package with the following command:

```
$ make test
```
