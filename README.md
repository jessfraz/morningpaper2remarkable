# morningpaper2remarkable

[![Travis CI](https://img.shields.io/travis/jessfraz/morningpaper2remarkable.svg?style=for-the-badge)](https://travis-ci.org/jessfraz/morningpaper2remarkable)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/jessfraz/morningpaper2remarkable)
[![Github All Releases](https://img.shields.io/github/downloads/jessfraz/morningpaper2remarkable/total.svg?style=for-the-badge)](https://github.com/jessfraz/morningpaper2remarkable/releases)

A bot to sync the [morning paper](https://blog.acolyer.org/) to a remarkable tablet.

This authenticates with your remarkable cloud account via the command line on
start. I hope to eventually make it run on my remarkable and not have to deal
with the cloud.

**Table of Contents**

<!-- toc -->

- [Installation](#installation)
    + [Binaries](#binaries)
    + [Via Go](#via-go)
    + [Running with Docker](#running-with-docker)
- [Usage](#usage)

<!-- tocstop -->

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/jessfraz/morningpaper2remarkable/releases).

#### Via Go

```console
$ go get github.com/jessfraz/morningpaper2remarkable
```

#### Running with Docker

**Authentication**

```console
$ touch ${HOME}/.rmapi

$ docker run --rm -it \
    --name morningpaper2remarkable \
    -v "${HOME}/.rmapi:/home/user/.rmapi:rw" \
    r.j3ss.co/morningpaper2remarkable --once

# Enter your one time auth code.
```

**Run it in daemon mode with our auth code**

```console
# You need to have already authed and have a .rmapi api file for this to 
# work in daemon mode.
$ docker run -d --restart always \
    --name morningpaper2remarkable \
    -v "${HOME}/.rmapi:/home/user/.rmapi:ro" \
    r.j3ss.co/morningpaper2remarkable --interval 20h
```

## Usage

```console
$ morningpaper2remarkable -h
morningpaper2remarkable -  A bot to sync the morning paper to a remarkable tablet.

Usage: morningpaper2remarkable <command>

Flags:

  -d, --debug  enable debug logging (default: false)
  --dir        directory to store the downloaded papers in (default: morningpaper)
  --interval   update interval (ex. 5ms, 10s, 1m, 3h) (default: 18h)
  --once       run once and exit, do not run as a daemon (default: false)

Commands:

  version  Show the version information.
```