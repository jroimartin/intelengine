# intelengine

## Introduction

intelengine aims to be an information gathering and exploitation architecture,
it is based on the use of transforms, that convert one data type into
another. For instance, a simple transform would be obtaining a list of
domains from an IP address or a location history from a twitter nickname.

## Main goals

The main goals of intelengine can be summarized in:

* Simplicity
* Modularity
* Scalability
* RESTful
* Programming language agnostic

## Architecture

intelengine consists in a client-server architecture.

intelsrv, the server component, is an HTTP server that exposes a REST API, that
allows the communication between server and clients. The mission of intelsrv is
handling command execution requests and transmit the output of the issued commands
to the client, as well as taking care of error handling, concurrency, caching,
server's OS abstraction, etc.

The client can be any program able to interact with the intelsrv's REST API.

## Commands

Commands are splitted in two parts:

* The definition file (cmd file)
* The implementation (standalone executable)

The command definition file is a JSON file that defines how the command is
called. It must include the following information:

* description: Description of the command functionality
* path: Path of the executable that will be called when the command is executed
* args: Arguments passed to the executable when it is called
* class: Command class

The following snippet shows a dummy cmd file:

```json
{
	"description": "echo request's body",
	"cmd": "cat",
	"args": [],
	"class": "debug"
}
```

Also, the definition files must have the extension ".cmd", being the name of the
command the name of the file without this extension.

The command implementation is an standalone executable that implements the
command functionality. By convention, it must wait for JSON input by STDIN and
write its output in JSON format to STDOUT. Also, it must exit with the return
value 0 when the execution finished correctly, or any other value on error.

The input of the command is the body of the PUT request sent to the intelsrv's
path "/cmd/exec/<command_name>". On the other hand, the output of the command
will be returned to the client in the response body if the command exited
successfully. Otherwise, if the command exited with error, an HTTP 500 error code
is returned to the client.

Due to these design principles, commands can be implemented in any programming
language that can read from STDIN and write to STDOUT.

## Transforms vs Commands

The word "command" was chosen rather than "transform", because a transform can be
considered as a particular class of command. intelengine is not only aimed at
being used for data gathering but also for exploitation, crawlering, etc.

## intelsrv's routes

The following routes are configured by default:

* GET /cmd/refresh: Refresh comand list
* GET /cmd/list: List supported command
* POST /cmd/exec/<cmdname>: Execute the command <cmdname>
