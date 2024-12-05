# openapi-gen

openapi-gen is a powerful tool for generating server-side and client-side code based on OpenAPI specifications.
It supports the following programming languages:

- Python
- Golang
- TypeScript

## Features:

- Multi-language Support: Easily generate server and client code for different programming languages.
- Configuration-driven: Customize the OpenAPI file path and output directory via a configuration file.
- Automated Workflow: Simplify the process from OpenAPI specification to code generation, reducing manual effort.

## Configfile

config example:

```.n-cl.yaml
golang:
  project: golang
  server: generated/server
  client: generated/client
  models: generated/models
  message: generated/message

typescript:
  project: typescript
  client: src/generated/client
  models: src/generated/models
  message: src/generated/message

python:
  project: python
  server: generated/server
  models: generated/models
  message: generated/message
```

## Command

```
Usage: main [<work-dir> [<config-path>]] [flags]

Arguments:
  [<work-dir>]       command work dir.
  [<config-path>]    config path, base on work dir.

Flags:
  -h, --help    Show context-sensitive help.
```

## Usage

Using openapi-gen involves a few simple steps:

### Step 1: Create a Configuration File

Prepare a configuration file `.n-cl.yaml` to specify the OpenAPI file path, target language, and output directory.

Step 2: Run the Command
Execute the following command to generate the desired code:

```
go run go run cmd/n-cl/main.go <work-dir>
```
