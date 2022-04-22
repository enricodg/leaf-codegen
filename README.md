# Leaf Code
This repository contains code generator for https://github.com/paulusrobin/leaf-utilities.

## Prerequisites

1. Go 1.18

## Installation

1. Clone this repository
```shell
$ git clone https://github.com/enricodg/leaf-codegen
```
2. Go to the directory, and install dependencies and binaries
```shell
$ cd leaf-codegen
$ go mod tidy
$ go install
```
3. You should be able to use the CLI by using `leaf-codegen`
```shell
$ leaf-codegen
NAME:
   Leaf Code Generator - Supporting leaf framework to initialize project

USAGE:
   leaf-codegen command [command options] [arguments...]

VERSION:
   v1.0.0

DESCRIPTION:
   CLI Leaf code generator

COMMANDS:
   init     init --project <project URL>
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)
```

## Usage

### Initialize Project

Create your directory / git repository and pass the value to init command, example:
```shell
$ leaf-codegen init --project github.com/enricodg/leaf-example
```

It will generate directories & files that looks like the following structure:
```
.
├── README.md
├── cmd
│   └── leaf-example
│       └── main.go
├── generateMock.sh
├── go.mod
├── go.sum
├── internal
│   ├── inbound
│   │   ├── di.go
│   │   └── http
│   │       ├── health
│   │       │   ├── check.go
│   │       │   ├── controller.go
│   │       │   └── routes.go
│   │       └── routes.go
│   ├── outbound
│   │   └── di.go
│   └── usecases
│       └── di.go
└── pkg
    ├── config
    │   ├── configApp.go
    │   ├── configNewRelic.go
    │   ├── configSentry.go
    │   └── di.go
    ├── di
    │   └── di.go
    └── resource
        ├── di.go
        ├── injection
        │   ├── logger.go
        │   ├── tracer.go
        │   ├── translator.go
        │   └── validator.go
        └── resource.go

13 directories, 23 files
```
Make sure you add all dependencies to the `go.sum` by running `go mod tidy`

## Example

For further example please visit https://github.com/enricodg/leaf-example.