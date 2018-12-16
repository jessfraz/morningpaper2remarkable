# morningpaper2remarkable

[![Travis CI](https://img.shields.io/travis/jessfraz/morningpaper2remarkable.svg?style=for-the-badge)](https://travis-ci.org/jessfraz/morningpaper2remarkable)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/jessfraz/morningpaper2remarkable)
[![Github All Releases](https://img.shields.io/github/downloads/jessfraz/morningpaper2remarkable/total.svg?style=for-the-badge)](https://github.com/jessfraz/morningpaper2remarkable/releases)

A bot to sync the morning paper to a remarkable tablet.

This authenticates with your remarkable cloud account via the command line on
start. I hope to eventually make it run on my remarkable and not have to deal
with the cloud.

* [Installation](README.md#installation)
   * [Binaries](README.md#binaries)
   * [Via Go](README.md#via-go)
* [Usage](README.md#usage)

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/jessfraz/morningpaper2remarkable/releases).

#### Via Go

```console
$ go get github.com/jessfraz/morningpaper2remarkable
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
