// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

intelworker is the worker component of intelengine.

It is responsible for executing the commands issued by clients and transmit
their results back to intelsrv. intelworker is designed to be programming
language agnostic, so commands can be coded using any language that can read
from STDIN and write to STDOUT.

Usage:
	intelworker config

The config file have the following format:
	{
		"Worker": {
			"CmdDir": "/path/to/cmds"
		},
		"Broker": {
			"URI": "amqp://amqp_broker:5672",
			"CertFile": "/path/to/cert.pem",
			"KeyFile": "/path/to/key.pem",
			"CAFile": "/path/to/cacert.pem",
			"Queue": "queue_name",
			"Exchange": "exchange_name"
		}
	}

Commands live with intelworker and are splitted in two parts:

	Definition file (cmd file)
	Implementation (standalone executable)

The command definition file is a JSON file that defines how the command is
called. It must include the following information:

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

Also, the definition files must have the extension ".cmd", being the name of the
command the name of the file without this extension.

The command's implementation is an standalone executable that implements the
command's functionality. By convention, it must wait for JSON input via STDIN
and write its output in JSON format to STDOUT. Also, it must exit with the
return value 0 when the execution finished correctly, or any other value on
error.

The input of the command is the body of the POST request sent to the intelsrv's
path "/cmd/exec/<cmdname>". When the users makes this request, an unique ID
will be generated and returned in the response. This way it is possible to
retrieve the result of the command sending a GET request to
"/cmd/result/<uuid>", being the output of the command returned to the client
in the response body. If the command exited with error, this error will be
returned in the "Error" field within the JSON response.

Commands must take care of the input and output types specified in their
definition file. Also, input and output must be treated as arrays of those
types. For instance, if the input type is "IP", the command should expect an
array of IPs as input.

Due to these design principles, commands can be implemented in any programming
language that can read from STDIN and write to STDOUT.

*/
package main
