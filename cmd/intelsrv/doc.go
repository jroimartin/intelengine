// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*

intelsrv is the server component of intelengine.

It is an HTTP server that exposes a REST API,
that allows the communication between server and clients. The mission of
intelsrv is handling execution flows and distribute tasks between the different
intelworker present in the architecture. Besides that, it also taskes care of
error handling, concurrency and caching

Usage:
	intelsrv config

The config file have the following format:
	{
		"Server": {
			"Addr": "0.0.0.0:8080"

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

The following default routes can be used by clients to control
intelsrv:

	GET /cmd/list: List supported commands
	POST /cmd/exec/<cmdname>: Execute the command <cmdname>
	GET /cmd/result/<uuid>: Retrieve the result of the command linked
		to the UUID <uuid>

*/
package main
