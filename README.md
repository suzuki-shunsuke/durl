# durl

[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](http://godoc.org/github.com/suzuki-shunsuke/durl)
[![Build Status](https://cloud.drone.io/api/badges/suzuki-shunsuke/durl/status.svg)](https://cloud.drone.io/suzuki-shunsuke/durl)
[![Docker Repository on Quay](https://quay.io/repository/suzuki_shunsuke/durl/status "Docker Repository on Quay")](https://quay.io/repository/suzuki_shunsuke/durl)
[![codecov](https://codecov.io/gh/suzuki-shunsuke/durl/branch/master/graph/badge.svg)](https://codecov.io/gh/suzuki-shunsuke/durl)
[![Go Report Card](https://goreportcard.com/badge/github.com/suzuki-shunsuke/durl)](https://goreportcard.com/report/github.com/suzuki-shunsuke/durl)
[![GitHub last commit](https://img.shields.io/github/last-commit/suzuki-shunsuke/durl.svg)](https://github.com/suzuki-shunsuke/durl)
[![GitHub tag](https://img.shields.io/github/tag/suzuki-shunsuke/durl.svg)](https://github.com/suzuki-shunsuke/durl/releases)
[![License](http://img.shields.io/badge/license-mit-blue.svg?style=flat-square)](https://raw.githubusercontent.com/suzuki-shunsuke/durl/master/LICENSE)

CLI tool to check whether dead urls are included in files.

* [Motivation](#motivation)
* [Install](#install)
* [Docker Image](#docker-image)
* [Getting Started](#getting-started)
* [Configuration](#configuration)
* [Change Log](https://github.com/suzuki-shunsuke/durl/releases)

## Motivation

Suppose you develop an oss and write documents about it.
You would want many users to use it.
But if it's documents include dead urls,
maybe users are disappointed and give up it even if it is good.

How sad it is!

So we have developed this tool.
It is good to use this tool at CI.

Of course, you can use durl other than oss documents.
For example, you can also check your blog posts with durl.

## Install

`durl` is written with Golang and binary is distributed at [release page](https://github.com/suzuki-shunsuke/durl/releases), so installation is easy.

## Docker Image

We provide busybox based docker image installed `durl`.

https://quay.io/repository/suzuki_shunsuke/durl

You can try to use durl without installation, and this is useful for CI.

```
$ docker run -ti --rm -v $PWD:/workspace -w /workspace quay.io/suzuki_shunsuke/durl sh
# echo foo.txt | durl check
```

## Getting Started

At first generate the configuration file.

```
# Generate .durl.yml
$ durl init
```

Generate a file included dead url.

```
$ cat << EOF > bar.txt
heredoc> foo
https://github.com/suzuki-shunsuke/durl
Please see https://github.com/suzuki-shunsuke/dead-repository .
EOF
```

Then check the file with `durl check`.
`durl check` accepts file paths as stdin.

```
$ echo foo.txt | durl check
[foo.txt] https://github.com/suzuki-shunsuke/dead-repository is dead (404)
```

It is good to use `durl` combining with the `find` command.

```
find . \
  -type d -name node_modules -prune -o \
  -type d -name .git -prune -o \
  -type d -name vendor -prune -o \
  -type f -print | \
  durl check || exit 1
```

## Ignore urls

* [check only urls whose scheme are "http" or "https"](https://github.com/suzuki-shunsuke/durl/issues/10)
* [ignore urls whose host matches the black list (ex. "localhost", "example.com")](https://github.com/suzuki-shunsuke/durl/issues/11)

## Configuration

```yaml
---
ignore_urls:
  - https://github.com/suzuki-shunsuke/ignore-repository
```

## Change Log

Please see [Releases](https://github.com/suzuki-shunsuke/durl/releases).

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) .

## License

[MIT](LICENSE)
