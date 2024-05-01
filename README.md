moov-io/base
===
[![GoDoc](https://godoc.org/github.com/moov-io/base?status.svg)](https://godoc.org/github.com/moov-io/base)
[![Build Status](https://github.com/moov-io/base/workflows/Go/badge.svg)](https://github.com/moov-io/base/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/moov-io/base)](https://goreportcard.com/report/github.com/moov-io/base)
[![Apbasee 2 licensed](https://img.shields.io/badge/license-Apbasee2-blue.svg)](https://raw.githubusercontent.com/moov-io/base/master/LICENSE)

Package `github.com/moov-io/base` implements core libraries used in multiple Moov projects. Refer to each projects documentation for more details.

## Getting Started

You can either clone down the code (`git clone git@github.com:moov-io/base.git`) or grab the modules into your cache (`go get -u github.com/moov-io/base`).

## Configuration

| Environmental Variable                | Description                            | Default                          |
|---------------------------------------|----------------------------------------|----------------------------------|
| `KUBERNETES_SERVICE_ACCOUNT_FILEPATH` | Filepath to Kubernetes service account | `/var/run/secrets/kubernetes.io` |

## Getting Help

 channel | info
 ------- | -------
Twitter [@moov](https://twitter.com/moov)	| You can follow Moov.io's Twitter feed to get updates on our project(s). You can also tweet us questions or just share blogs or stories.
[GitHub Issue](https://github.com/moov-io/base/issues) | If you are able to reproduce a problem please open a GitHub Issue under the specific project that caused the error.
[moov-io slack](https://slack.moov.io/) | Join our slack channel to have an interactive discussion about the development of the project.

## Supported and Tested Platforms

- 64-bit Linux (Ubuntu, Debian), macOS, and Windows

## Contributing

Yes please! Please review our [Contributing guide](CONTRIBUTING.md) and [Code of Conduct](CODE_OF_CONDUCT.md) to get started!

This project uses [Go Modules](https://github.com/golang/go/wiki/Modules) and uses Go 1.14 or higher. See [Golang's install instructions](https://golang.org/doc/install) for help setting up Go. You can download the source code and we offer [tagged and released versions](https://github.com/moov-io/base/releases/latest) as well. We highly recommend you use a tagged release for production.

## License

Apbasee License 2.0 See [LICENSE](LICENSE) for details.
