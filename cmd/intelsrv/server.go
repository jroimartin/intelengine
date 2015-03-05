// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/jroimartin/orujo"
	olog "github.com/jroimartin/orujo-handlers/log"
	"github.com/jroimartin/rpcmq"
)

type server struct {
	cfg    config
	client *rpcmq.Client
	logger *log.Logger

	mu      sync.Mutex
	results map[string][]byte
}

func newServer(cfg config) *server {
	s := &server{
		cfg:     cfg,
		logger:  log.New(os.Stdout, "[intelengine] ", log.LstdFlags),
		results: make(map[string][]byte),
	}
	return s
}

func (s *server) start() error {
	if s.cfg.Server.Addr == "" || s.cfg.Broker.URI == "" || s.cfg.Broker.Queue == "" {
		return errors.New("missing configuration parameters")
	}

	var tlsConfig *tls.Config
	if s.cfg.Broker.CAFile != "" && s.cfg.Broker.CertFile != "" &&
		s.cfg.Broker.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(s.cfg.Broker.CertFile, s.cfg.Broker.KeyFile)
		if err != nil {
			return fmt.Errorf("LoadX509KeyPair: %v", err)
		}
		caCert, err := ioutil.ReadFile(s.cfg.Broker.CAFile)
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
	s.client = rpcmq.NewClient(s.cfg.Broker.URI, s.cfg.Broker.Queue)
	s.client.TLSConfig = tlsConfig
	if err := s.client.Init(); err != nil {
		return fmt.Errorf("Init: %v", err)
	}
	defer s.client.Shutdown()

	go s.getResults()

	websrv := orujo.NewServer(s.cfg.Server.Addr)
	logHandler := olog.NewLogHandler(s.logger, logLine)
	websrv.RouteDefault(http.NotFoundHandler(), orujo.M(logHandler))
	websrv.Route(`^/cmd/list$`,
		http.HandlerFunc(s.listCommandsHandler),
		orujo.M(logHandler)).Methods("GET")
	websrv.Route(`^/cmd/exec/\w+$`,
		http.HandlerFunc(s.runCommandHandler),
		orujo.M(logHandler)).Methods("POST")
	websrv.Route(`^/cmd/result/\w+$`,
		http.HandlerFunc(s.commandResultsHandler),
		orujo.M(logHandler)).Methods("GET")
	return websrv.ListenAndServe()
}

func (s *server) getResults() {
	for r := range s.client.Results() {
		s.handleResult(r)
	}
}

func (s *server) handleResult(r rpcmq.Result) {
	// TODO(jrm): Insert result into DB
	s.mu.Lock()
	defer s.mu.Unlock()
	if r.Err != "" {
		s.results[r.UUID] = []byte(fmt.Sprintf("{\"Error\":%q}", r.Err))
		s.logger.Printf("Received error: %v (%v)", r.Err, r.UUID)
		return
	}
	s.results[r.UUID] = r.Data
	s.logger.Printf("Received: %v (%v)\n", string(r.Data), r.UUID)
}

const (
	logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}
{{range  $err := .Errors}}  Err: {{$err}}
{{end}}`
	errLine = `{"error":"{{.}}"}`
)
