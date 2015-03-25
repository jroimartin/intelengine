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
	Worker workerConfig
	Broker brokerConfig
}

type workerConfig struct {
	CmdDir string
}

type brokerConfig struct {
	URI       string
	CertFile  string
	KeyFile   string
	CAFile    string
	MsgsQueue string
	Exchange  string
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

	w := newWorker(cfg)
	log.Fatalln(w.start())
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: intelworker config")
}
