# argodocs

`argodocs` is a tool for generating high-quality documentation from argo workflow templates.

## Usage
```shell
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
```shell
argodocs generate **/*.yaml --output-prefix=../docs/
```