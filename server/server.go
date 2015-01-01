// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/howeyc/fsnotify"
	"github.com/jroimartin/orujo"
	olog "github.com/jroimartin/orujo-handlers/log"
)

type Server struct {
	Logger   *log.Logger
	Addr     string
	CmdDir   string
	commands []command
}

func NewServer() *Server {
	s := new(Server)
	s.Logger = log.New(os.Stdout, "[intelengine] ", log.LstdFlags)
	return s
}

func (s *Server) Start() error {
	if err := s.setupWatcher(); err != nil {
		return err
	}

	if err := s.setupServer(); err != nil {
		return err
	}

	return nil
}

func (s *Server) setupServer() error {
	if s.Addr == "" {
		return errors.New("Server.Addr cannot be an empty string")
	}

	os := orujo.NewServer(s.Addr)

	logHandler := olog.NewLogHandler(s.Logger, logLine)

	// TODO: Add routes
	os.RouteDefault(http.NotFoundHandler(), orujo.M(logHandler))

	if err := os.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) setupWatcher() error {
	if s.CmdDir == "" {
		return errors.New("Server.CmdDir cannot be an empty string")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go s.trackCmds(watcher)

	if err = watcher.Watch(s.CmdDir); err != nil {
		return err
	}

	return nil
}

func (s *Server) trackCmds(watcher *fsnotify.Watcher) {
	for {
		select {
		case ev := <-watcher.Event:
			if path.Ext(ev.Name) != ".cmd" {
				continue
			}
			s.Logger.Print("watcher evernt:", ev)
		case err := <-watcher.Error:
			s.Logger.Print("watcher error:", err)
		}
	}
}

const logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}
{{range  $err := .Errors}}  Err: {{$err}}
{{end}}`
