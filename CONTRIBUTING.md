# Contributing

## Check before send a pull request

* Commit message format conforms the [AngularJS Commit Message Format](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#commits)
* [Commit message type is appropriate](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#type)
* [If the pull request includes breaking changes, please describe them](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#footer)
* Document and code comment is updated

## Requirements

* [npm](https://www.npmjs.com/)
* [Golang](https://golang.org/)
* [golangci-lint](https://github.com/golangci/golangci-lint)
* [goreleaser](https://goreleaser.com/)
* [gomic](https://github.com/suzuki-shunsuke/gomic)
* [drone CLI](https://docs.drone.io/cli/install/)

We use node libraries and npm scripts for development.
Please see [package.json](https://github.com/suzuki-shunsuke/durl/blob/master/package.json) .

## Set up

```
$ npm i
$ go mod download
```

## Lint

```
# Lint with go vet.
$ npm run vet
# Lint with golangci-lint
$ npm run lint
```

## Format codes with gofmt

```
$ npm run fmt
```

## Test

```
$ npm t
# Test with drone ci
# https://docs.drone.io/cli/drone-exec/
$ npm run ci-local
```

## Generate mocks

```
# Generate mocks for tests of durl.
$ npm run gen-mock
```

## Commit Message Format

The commit message format of this project conforms to the [AngularJS Commit Message Format](https://github.com/angular/angular.js/blob/master/DEVELOPERS.md#commits).
We validate the commit message with git's `commit-msg` hook using [commitlint](http://marionebl.github.io/commitlint/#/) and [husky](https://www.npmjs.com/package/husky), so you have to install them before commit.

```
$ npm i
```

## Release

```
$ npm run tag <tag>
$ git push
$ git push --tags
```

Tag must start with "v".
`npm run tag` command updates [internal/domain/version.go](https://github.com/suzuki-shunsuke/durl/blob/master/internal/domain/version.go) and commit and creates a tag.
When we push a tag to GitHub, ci is run and durl is built and uploaded to [GitHub Relases](https://github.com/suzuki-shunsuke/durl/releases) .

## CI

We use [drone cloud](https://cloud.drone.io/suzuki-shunsuke/durl).

Please see [.drone.yml](https://github.com/suzuki-shunsuke/durl/blob/master/.drone.yml) .
