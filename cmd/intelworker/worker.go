package main

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/jroimartin/rpcmq"
)

type worker struct {
	cfg    config
	server *rpcmq.Server
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

	select {}
}
