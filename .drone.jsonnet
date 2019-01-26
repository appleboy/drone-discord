local name = 'drone-discord';

local PipelineTesting = {
  kind: "pipeline",
  name: "testing",
  platform: {
    os: "linux",
    arch: "amd64",
  },
  steps: [
    {
      name: "vet",
      image: "golang:1.11",
      pull: "always",
      environment: {
        GO111MODULE: "on",
      },
      commands: [
        "make vet",
      ],
      volumes: [
        {
          name: "gopath",
          path: "/go",
        },
      ],
    },
    {
      name: "lint",
      image: "golang:1.11",
      pull: "always",
      environment: {
        GO111MODULE: "on",
      },
      commands: [
        "make lint",
      ],
      volumes: [
        {
          name: "gopath",
          path: "/go",
        },
      ],
    },
    {
      name: "misspell",
      image: "golang:1.11",
      pull: "always",
      environment: {
        GO111MODULE: "on",
      },
      commands: [
        "make misspell-check",
      ],
      volumes: [
        {
          name: "gopath",
          path: "/go",
        },
      ],
    },
    {
      name: "test",
      image: "golang:1.11",
      pull: "always",
      environment: {
        GO111MODULE: "on",
        WEBHOOK_ID: { "from_secret": "webhook_id" },
        WEBHOOK_TOKEN: { "from_secret": "webhook_token" },
      },
      commands: [
        "make test",
        "make coverage",
      ],
      volumes: [
        {
          name: "gopath",
          path: "/go",
        },
      ],
    },
    {
      name: "codecov",
      image: "robertstettner/drone-codecov",
      pull: "always",
      settings: {
        token: { "from_secret": "codecov_token" },
      },
    },
  ],
  volumes: [
    {
      name: "gopath",
      temp: {},
    },
  ],
  trigger: {
    branch: [ "master" ],
  },
};

local PipelineBuild(name, os="linux", arch="amd64") = {
  kind: "pipeline",
  name: os + "-" + arch,
  platform: {
    os: os,
    arch: arch,
  },
  steps: [
    {
      name: "build-push",
      image: "golang:1.11",
      pull: "always",
      environment: {
        CGO_ENABLED: "0",
        GO111MODULE: "on",
      },
      commands: [
        "go build -v -ldflags \"-X main.build=${DRONE_BUILD_NUMBER}\" -a -o release/" + os + "/" + arch + "/" + name,
      ],
      when: {
        event: [ "push", "pull_request" ],
      },
    },
    {
      name: "build-tag",
      image: "golang:1.11",
      pull: "always",
      environment: {
        CGO_ENABLED: "0",
        GO111MODULE: "on",
      },
      commands: [
        "go build -v -ldflags \"-X main.version=${DRONE_TAG##v} -X main.build=${DRONE_BUILD_NUMBER}\" -a -o release/" + os + "/" + arch + "/" + name,
      ],
      when: {
        event: [ "tag" ],
      },
    },
    {
      name: "executable",
      image: "golang:1.11",
      pull: "always",
      commands: [
        "./release/" + os + "/" + arch + "/" + name + " --help",
      ],
    },
    {
      name: "dryrun",
      image: "plugins/docker:" + os + "-" + arch,
      pull: "always",
      settings: {
        dry_run: true,
        tags: os + "-" + arch,
        dockerfile: "docker/Dockerfile." + os + "." + arch,
        repo: "appleboy/" + name,
        username: { "from_secret": "docker_username" },
        password: { "from_secret": "docker_password" },
      },
      when: {
        event: [ "pull_request" ],
      },
    },
    {
      name: "publish",
      image: "plugins/docker:" + os + "-" + arch,
      pull: "always",
      settings: {
        auto_tag: true,
        auto_tag_suffix: os + "-" + arch,
        dockerfile: "docker/Dockerfile." + os + "." + arch,
        repo: "appleboy/" + name,
        username: { "from_secret": "docker_username" },
        password: { "from_secret": "docker_password" },
      },
      when: {
        event: [ "push", "tag" ],
      },
    },
  ],
  depends_on: [
    "testing",
  ],
  trigger: {
    branch: [ "master" ],
  },
};

local PipelineNotifications = {
  kind: "pipeline",
  name: "notifications",
  platform: {
    os: "linux",
    arch: "amd64",
  },
  clone: {
    disable: true,
  },
  steps: [
    {
      name: "microbadger",
      image: "plugins/webhook:1",
      pull: "always",
      settings: {
        url: { "from_secret": "microbadger_url" },
      },
    },
  ],
  depends_on: [
    "linux-amd64",
    "linux-arm64",
    "linux-arm",
  ],
  trigger: {
    branch: [ "master" ],
    event: [ "push", "tag" ],
  },
};

[
  PipelineTesting,
  PipelineBuild(name, "linux", "amd64"),
  PipelineBuild(name, "linux", "arm64"),
  PipelineBuild(name, "linux", "arm"),
  PipelineNotifications,
]
