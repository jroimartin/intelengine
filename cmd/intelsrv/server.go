// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"

	"github.com/jroimartin/orujo"
	olog "github.com/jroimartin/orujo-handlers/log"
)

type server struct {
	addr     string
	cmdDir   string
	logger   *log.Logger
	commands map[string]*command
	mutex    sync.RWMutex
}

func newServer(addr, cmdDir string) *server {
	s := &server{
		addr:   addr,
		cmdDir: cmdDir,
		logger: log.New(os.Stdout, "[intelengine] ", log.LstdFlags),
	}
	return s
}

func (s *server) start() error {
	if s.addr == "" || s.cmdDir == "" {
		return errors.New("server.addr and server.cmdDir cannot be empty strings")
	}

	s.refreshCommands()

	websrv := orujo.NewServer(s.addr)

	logHandler := olog.NewLogHandler(s.logger, logLine)

	websrv.RouteDefault(http.NotFoundHandler(), orujo.M(logHandler))

	websrv.Route(`^/cmd/refresh$`,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.refreshCommands()
		}),
		http.HandlerFunc(s.listCommandsHandler),
		orujo.M(logHandler)).Methods("GET")

	websrv.Route(`^/cmd/list$`,
		http.HandlerFunc(s.listCommandsHandler),
		orujo.M(logHandler)).Methods("GET")

	websrv.Route(`^/cmd/exec/\w+$`,
		http.HandlerFunc(s.runCommandHandler),
		orujo.M(logHandler)).Methods("POST")

	return websrv.ListenAndServe()
}

func (s *server) refreshCommands() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.commands = make(map[string]*command)

	files, err := ioutil.ReadDir(s.cmdDir)
	if err != nil {
		s.logger.Println("refreshCommands warning:", err)
		return
	}

	for _, f := range files {
		if f.IsDir() || path.Ext(f.Name()) != cmdExt {
			continue
		}

		filename := path.Join(s.cmdDir, f.Name())
		cmd, err := newCommand(filename)
		if err != nil {
			s.logger.Println("refreshCommands warning:", err)
			return
		}

		s.commands[cmd.Name] = cmd
		s.logger.Println("command registered:", cmd.Name)
	}
}

func (s *server) command(name string) *command {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, cmd := range s.commands {
		if cmd.Name == name {
			return cmd
		}
	}
	return nil
}

const logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}
{{range  $err := .Errors}}  Err: {{$err}}
{{end}}`
