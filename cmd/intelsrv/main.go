// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

type config struct {
	Server serverConfig
	Broker brokerConfig
}

type serverConfig struct {
	Addr string
}

type brokerConfig struct {
	URI      string
	CertFile string
	KeyFile  string
	CAFile   string
	Queue    string
	Exchange string
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	configFile := flag.Arg(0)
	f, err := os.Open(configFile)
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()

	var cfg config
	dec := json.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		log.Fatalln(err)
	}

	s := newServer(cfg)
	log.Fatalln(s.start())
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: intelsrv config")
}
