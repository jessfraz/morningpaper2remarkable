# morningpaper2remarkable

[![make-all](https://github.com/pseudo-su/morningpaper2remarkable/workflows/make%20all/badge.svg)](https://github.com/pseudo-su/morningpaper2remarkable/actions?query=workflow%3A%22make+all%22)
[![make-image](https://github.com/pseudo-su/morningpaper2remarkable/workflows/make%20image/badge.svg)](https://github.com/pseudo-su/morningpaper2remarkable/actions?query=workflow%3A%22make+image%22)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/pseudo-su/morningpaper2remarkable)
[![Github All Releases](https://img.shields.io/github/downloads/pseudo-su/morningpaper2remarkable/total.svg?style=for-the-badge)](https://github.com/pseudo-su/morningpaper2remarkable/releases)

A bot to sync the [morning paper](https://blog.acolyer.org/) to a remarkable tablet.

This authenticates with your remarkable cloud account via the command line on
start. I hope to eventually make it run on my remarkable and not have to deal
with the cloud.

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [morningpaper2remarkable](#morningpaper2remarkable)
  - [Installation](#installation)
      - [Binaries](#binaries)
      - [Via Go](#via-go)
      - [Running with Docker](#running-with-docker)
  - [Usage](#usage)
    - [Hidden Command](#hidden-command)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/pseudo-su/morningpaper2remarkable/releases).

#### Via Go

```console
$ go get github.com/pseudo-su/morningpaper2remarkable
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
    -v "${HOME}/.rmapi:/home/user/.rmapi:rw" \
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
  --pages      number of pages of papers to download (default: 1)

Commands:

  version  Show the version information.
```

### Hidden Command

I use the bot on my server but sometimes I just want a way to get a paper from
a URL to my remarkable from the command line.

I added a hidden command for that `download`.

You can use it like the following:

```bash
$ morningpaper2remarkable download http://nickm.com/trope_tank/10_PRINT_121114.pdf "10 PRINT"
```

This will download the PDF from the URL at `arg[0]` put it in a folder, default
named `papers` and name the PDF in that folder `arg[1]`, which above is `"10
PRINT"`.

You can change the folder name with the `--dataDir` flag.
