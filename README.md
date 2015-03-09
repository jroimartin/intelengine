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

**intelsrv**, the server component, is an HTTP server that exposes a REST API,
that allows the communication between server and clients. The mission of
intelsrv is handling execution flows and distribute tasks between the different
intelworker present in the architecture. Besides that, it also taskes care of
error handling, concurrency and caching

**intelworker**, the worker component is responsible for executing the commands
issued by clients and transmit their results back to intelsrv. intelworker is
designed to be programming language agnostic, so commands can be coded using
any language that can read from STDIN and write to STDOUT.

Finally, the **client** can be any program able to interact with the intelsrv's
REST API.

It is important to note that the communication between the different instances
of intelserver and intelworker is carried out via a message broker using the
amqp protocol.

```
+--------+  http    +------------+  amqp    +--------+  amqp    +---------------+
| client |----+---->| intelsrv_1 |----+---->| BROKER |----+---->| intelworker_1 |
+--------+    |     +------------+    |     +--------+    |     +---------------+
              |     +------------+    |                   |     +---------------+
              +---->| intelsrv_2 |----+                   +---->| intelworker_2 |
              |     +------------+    |                   |     +---------------+
              |     +------------+    |                   |     +---------------+
              +---->| intelsrv_n |----+                   +---->| intelworker_m |
                    +------------+                              +---------------+
```

## Commands

Commands live with intelworker and are splitted in two parts:

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

The input of the command is the body of the POST request sent to the intelsrv's
path "/cmd/exec/\<cmdname\>". When the users makes this request, an unique ID
will be generated and returned in the response. This way it is possible to
retrieve the result of the command sending a GET request to
"/cmd/result/\<uuid\>", being the output of the command returned to the client
in the response body. If the command exited with error, this error will be
returned in the "Error" field within the JSON response.

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

* **GET /cmd/list**: List supported commands
* **POST /cmd/exec/\<cmdname\>**: Execute the command \<cmdname\>
* **GET /cmd/result/\<uuid\>**: Retrieve the result of the command linked
	to the UUID \<uuid\>
