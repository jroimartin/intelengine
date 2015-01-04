// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

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

type Server struct {
	Addr   string
	CmdDir string

	logger   *log.Logger
	commands map[string]*command
	mutex    sync.RWMutex
}

func NewServer() *Server {
	s := new(Server)
	s.logger = log.New(os.Stdout, "[intelengine] ", log.LstdFlags)
	return s
}

func (s *Server) Start() error {
	if s.Addr == "" || s.CmdDir == "" {
		return errors.New("Server.Addr and Server.CmdDir cannot be empty strings")
	}

	s.refreshCommands()

	websrv := orujo.NewServer(s.Addr)

	logHandler := olog.NewLogHandler(s.logger, logLine)

	websrv.RouteDefault(http.NotFoundHandler(), orujo.M(logHandler))

	websrv.Route(`^/cmd/refresh$`,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			s.refreshCommands()
		}),
		http.HandlerFunc(s.listCommandsHandler),
		orujo.M(logHandler))

	websrv.Route(`^/cmd/list$`,
		http.HandlerFunc(s.listCommandsHandler),
		orujo.M(logHandler))

	websrv.Route(`^/cmd/exec/\w+$`,
		http.HandlerFunc(s.runCommandHandler),
		orujo.M(logHandler))

	return websrv.ListenAndServe()
}

func (s *Server) refreshCommands() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.commands = make(map[string]*command)

	files, err := ioutil.ReadDir(s.CmdDir)
	if err != nil {
		s.logger.Println("refreshCommands warning:", err)
		return
	}

	for _, f := range files {
		if f.IsDir() || path.Ext(f.Name()) != cmdExt {
			continue
		}

		filename := path.Join(s.CmdDir, f.Name())
		cmd, err := newCommand(filename)
		if err != nil {
			s.logger.Println("refreshCommands warning:", err)
			return
		}

		s.commands[cmd.Name] = cmd
		s.logger.Println("command registered:", cmd.Name)
	}
}

func (s *Server) command(name string) *command {
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
