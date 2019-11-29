# Contributing

## Check before send a pull request

* [If the pull request includes breaking changes, please describe them](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#footer)
* Document and code comment is updated

## Requirements

* [cmdx](https://github.com/suzuki-shunsuke/cmdx)
* [Golang](https://golang.org/)
* [golangci-lint](https://github.com/golangci/golangci-lint)
* [goreleaser](https://goreleaser.com/)
* [gomic](https://github.com/suzuki-shunsuke/gomic)
* [drone CLI](https://docs.drone.io/cli/install/)

We use [cmdx](https://github.com/suzuki-shunsuke/cmdx) for task runner.
You can check tasks by `cmdx -l`

```
$ cmdx -l
```

## Commit Message Format

The commit message format of this project conforms to the [AngularJS Commit Message Format](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#commits).

## CI

We use [drone cloud](https://cloud.drone.io/suzuki-shunsuke/durl).

Please see [.drone.yml](https://github.com/suzuki-shunsuke/durl/blob/master/.drone.yml) .
