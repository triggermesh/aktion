[![Go Report Card](https://goreportcard.com/badge/github.com/triggermesh/aktion)](https://goreportcard.com/report/github.com/triggermesh/aktion) [![CircleCI](https://circleci.com/gh/triggermesh/aktion/tree/master.svg?style=shield)](https://circleci.com/gh/triggermesh/aktion/tree/master)

A CLI for submitting [Github Actions](https://developer.github.com/actions/creating-workflows/workflow-configuration-options/#workflow-blocks) to [Knative Build](https://github.com/knative/build)

## Installation

With a working [Golang](https://golang.org/doc/install) environment do:

```
go get github.com/triggermesh/aktion
```

## Usage

```
aktion pipeline -f samples/main.workflow
```

## Support

We would love your feedback on this tool so don't hesitate to let us know what is wrong and how we could improve it, just file an [issue](https://github.com/triggermesh/aktion/issues/new)

## Code of Conduct

This plugin is by no means part of [CNCF](https://www.cncf.io/) but we abide by its [code of conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md)
