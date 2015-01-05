// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
)

var (
	addr   = flag.String("addr", ":8080", "HTTP service address")
	cmddir = flag.String("cmddir", "cmds", "directory containing command descriptions")
)

func main() {
	flag.Parse()
	s := newServer(*addr, *cmddir)
	log.Fatalln(s.start())
}
