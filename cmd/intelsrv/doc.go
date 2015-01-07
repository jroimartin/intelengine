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

	description: Description of the command functionality
	path: Path of the executable that will be called when the command is
		executed
	args: Arguments passed to the executable when it is called
	class: Command class

The following snippet shows a dummy cmd file:

	{
		"description": "echo request's body",
		"cmd": "cat",
		"args": [],
		"class": "debug"
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

Due to these design principles, commands can be implemented in any
programming language that can read from STDIN and write to STDOUT.

The following default routes can be used by clients to control
intelsrv:

	GET /cmd/refresh: Refresh comand list
	GET /cmd/list: List supported commands
	POST /cmd/exec/<cmdname>: Execute the command <cmdname>

*/
package main
