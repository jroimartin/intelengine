// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"

	"github.com/jroimartin/rpcmq"
)

const (
	sep = '|'
)

type worker struct {
	cfg    config
	server *rpcmq.Server

	mu       sync.RWMutex
	commands map[string]*command
}

func newWorker(cfg config) *worker {
	return &worker{cfg: cfg}
}

func (w *worker) start() error {
	if w.cfg.Worker.CmdDir == "" || w.cfg.Broker.URI == "" || w.cfg.Broker.Queue == "" {
		return errors.New("missing configuration parameters")
	}

	var tlsConfig *tls.Config
	if w.cfg.Broker.CAFile != "" && w.cfg.Broker.CertFile != "" &&
		w.cfg.Broker.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(w.cfg.Broker.CertFile, w.cfg.Broker.KeyFile)
		if err != nil {
			return fmt.Errorf("LoadX509KeyPair: %v", err)
		}
		caCert, err := ioutil.ReadFile(w.cfg.Broker.CAFile)
		if err != nil {
			return fmt.Errorf("ReadFile: %v", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      caCertPool,
		}
	}
	w.server = rpcmq.NewServer(w.cfg.Broker.URI, w.cfg.Broker.Queue)
	w.server.TLSConfig = tlsConfig
	if err := w.server.Init(); err != nil {
		return fmt.Errorf("Init: %v", err)
	}
	defer w.server.Shutdown()

	w.refreshCommands()
	if err := w.server.Register("listCommands", w.listCommands); err != nil {
		return err
	}
	if err := w.server.Register("execCommand", w.execCommand); err != nil {
		return err
	}

	select {}
}

func (w *worker) listCommands(data []byte) ([]byte, error) {
	if err := w.refreshCommands(); err != nil {
		return nil, fmt.Errorf("cannot refresh commands: %v", err)
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	b, err := json.Marshal(w.commands)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal commands: %v", err)
	}
	return b, nil
}

func (w *worker) refreshCommands() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.commands = map[string]*command{}
	if err := filepath.Walk(w.cfg.Worker.CmdDir, w.handleFile); err != nil {
		return err
	}
	return nil
}

func (w *worker) handleFile(filepath string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() || path.Ext(info.Name()) != cmdExt {
		return nil
	}

	cmd, err := newCommand(filepath)
	if err != nil {
		log.Printf("handleFile warning (%v): %v\n", info.Name(), err)
		return nil
	}

	w.commands[cmd.Name] = cmd
	log.Println("command registered:", cmd.Name)
	return nil
}

func (w *worker) execCommand(data []byte) ([]byte, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	sepIdx := bytes.IndexByte(data, sep)
	if sepIdx < 0 {
		return nil, errors.New("separator not found")
	}
	name := string(data[:sepIdx])
	br := bytes.NewReader(data[sepIdx+1:])

	cmd := w.command(name)
	if cmd == nil {
		return nil, fmt.Errorf("command not found: %v", name)
	}

	out, err := cmd.exec(br)
	if err != nil {
		return nil, fmt.Errorf("command execution error: %v", err)
	}

	return out, nil
}

func (w *worker) command(name string) *command {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if cmd, ok := w.commands[name]; ok {
		return cmd
	}
	return nil
}
