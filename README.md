# argodocs
[![Go Reference](https://pkg.go.dev/badge/pkg.go.dev/github.com/rohankmr414/argodocs.svg)](https://pkg.go.dev/github.com/rohankmr414/argodocs) [![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/rohankmr414/argodocs)](https://github.com/rohankmr414/argodocs/releases/tag/latest) [![build](https://github.com/rohankmr414/argodocs/actions/workflows/build.yaml/badge.svg)](https://github.com/rohankmr414/argodocs/actions/workflows/build.yaml) [![License: MIT](https://img.shields.io/badge/License-MIT-black.svg)](https://opensource.org/licenses/MIT)


`argodocs` is a tool to generate reference documentation for argo workflow templates.

## Installation


```
go install github.com/rohankmr414/argodocs
```
Add `$GOPATH/bin` to your `$PATH` or copy `$GOPATH/bin/argodocs` to your `$PATH`.
## Usage
```
Usage:
  argodocs [command]

Available Commands:
  generate    Generate docs from workflow manifest.
  help        Help about any command

Flags:
  -h, --help   help for argodocs

Use "argodocs [command] --help" for more information about a command.
```

## Generate
```
argodocs generate **/*.yaml --output-prefix=../docs/
```

## Contributing

If you're interested in contributing to argodocs, check out [CONTRIBUTING.md](CONTRIBUTING.md)

## License

Copyright (c) **Rohan Kumar** and **Junaid Rahim**. All rights reserved. Released under the [MIT](LICENSE) License