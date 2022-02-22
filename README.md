# argodocs

`argodocs` is a tool to generate reference documentation for argo workflow templates.

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
