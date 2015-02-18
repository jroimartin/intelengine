// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

intelsrv is the server component of intelengine.

It is an HTTP server that exposes a REST API, that allows the
communication between server and clients. The mission of intelsrv
is handling command execution requests and transmit the output of
the issued commands to the client, as well as taking care of error
handling, concurrency, caching, server's OS abstraction, etc.

Usage:
	intelsrv [flag]

The flags are:
	-addr=":8080": HTTP service address
	-cmddir="cmds": directory containing command definitions

Commands are splitted in two parts:

	Definition file (cmd file)
	Implementation (standalone executable)

The command definition file is a JSON file that defines how the
command is called. It must include the following information:

	Description: Description of the command functionality
	Path: Path of the executable that will be called when the command is
		executed
	Args: Arguments passed to the executable when it is called
	Input: Type of the input data
	Output: Type of the output data
	Parameters: Structure describing the type of the accepted parameters
	Group: Command category

The following snippet shows a dummy cmd file:

	{
		"Description": "echo request's body",
		"Cmd": "cat",
		"Args": [],
		"Input": "",
		"Output": "",
		"Parameters": "",
		"Group": "debug"
	}

Also, the definition files must have the extension ".cmd", being the
name of the command the name of the file without this extension.

The command's implementation is an standalone executable that
implements the command's functionality. By convention, it must wait
for JSON input via STDIN and write its output in JSON format to
STDOUT. Also, it must exit with the return value 0 when the execution
finished correctly, or any other value on error.

The input of the command is the body of the PUT request sent to the
intelsrv's path "/cmd/exec/<cmdname>". On the other hand, the
output of the command will be returned to the client in the response
body if the command exited successfully. Otherwise, if the command
exited with error, an HTTP 500 error code is returned to the client.

Commands must take care of the input and output types specified in their
definition file. Also, input and output must be treated as arrays of those
types. For instance, if the input type is "IP", the command should expect an
array of IPs as input.

Due to these design principles, commands can be implemented in any
programming language that can read from STDIN and write to STDOUT.

The following default routes can be used by clients to control
intelsrv:

	GET /cmd/refresh: Refresh command list
	GET /cmd/list: List supported commands
	POST /cmd/exec/<cmdname>: Execute the command <cmdname>

*/
package main
