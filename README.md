# durl

cli tool to check whether broken urls are included in files

## Motivation

Suppose you develop an oss and write documents about it.
You would want many users to use it.
But if it's documents include dead urls,
maybe users are disappointed and give up it even if it is good.

How sad it is!

So we have developed this tool.
It is good to use this tool at CI.

## Install

durl is written with Golang and binary is distributed at [release page](https://github.com/suzuki-shunsuke/durl/releases), so installation is easy.

## Getting Started

durl accepts file paths as stdin.

```
$ cat foo.txt
foo
http://example.com
Please see http://example.com/bar .

$ echo "foo.txt" | durl check
http://example.com/bar is broken (404)
```

## Configuration

```yaml
ignore_urls:
- http://example.com/bar
```

## Change Log

Please see [Releases](https://github.com/suzuki-shunsuke/durl/releases).

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) .

## License

[MIT](LICENSE)
