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

intelsrv, the **server** component, is an HTTP server that exposes a REST API,
	that allows the communication between server and clients. The mission of
	intelsrv is handling command execution requests and transmit the output of
	the issued commands to the client, as well as taking care of error handling,
	concurrency, caching, server's OS abstraction, etc.

The **client** can be any program able to interact with the intelsrv's REST API.

## Commands

Commands are splitted in two parts:

* **Definition file** (cmd file)
* **Implementation** (standalone executable)

The command **definition file** is a JSON file that defines how the command is
called. It must include the following information:

* **Description**: Description of the command functionality
* **Path**: Path of the executable that will be called when the command is executed
* **Args**: Arguments passed to the executable when it is called
* **Input**: Type of the input data
* **Output**: Type of the output data
* **Parameters**: Structure describing the type of the accepted parameters
* **Group**: Command category

The following snippet shows a dummy cmd file:

```json
{
	"Description": "echo request's body",
	"Cmd": "cat",
	"Args": [],
	"Input": "",
	"Output": "",
	"Parameters": "",
	"Group": "debug"
}
```

Also, the definition files must have the extension ".cmd", being the name of the
command the name of the file without this extension.

The command's **implementation** is an standalone executable that implements the
command's functionality. By convention, it must wait for JSON input via STDIN
and write its output in JSON format to STDOUT. Also, it must exit with the
return value 0 when the execution finished correctly, or any other value on
error.

The input of the command is the body of the PUT request sent to the intelsrv's
path "/cmd/exec/\<cmdname\>". On the other hand, the output of the command will
be returned to the client in the response body if the command exited
successfully. Otherwise, if the command exited with error, an HTTP 500 error
code is returned to the client.

Commands must take care of the input and output types specified in their
definition file. Also, input and output must be treated as arrays of those
types. For instance, if the input type is "IP", the command should expect an
array of IPs as input.

Due to these design principles, commands can be implemented in any programming
language that can read from STDIN and write to STDOUT.

## Transforms vs Commands

The word **command** was chosen rather than **transform**, because a transform
can be considered as a particular class of command. It's important to take into
account that intelengine is not only aimed at being used for data gathering
but also for exploitation, crawlering, etc.

## intelsrv's routes

The following routes are configured by default:

* **GET /cmd/refresh**: Refresh command list
* **GET /cmd/list**: List supported commands
* **POST /cmd/exec/\<cmdname\>**: Execute the command \<cmdname\>
