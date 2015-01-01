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

	"github.com/howeyc/fsnotify"
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
	s.initCommands()

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

	websrv := orujo.NewServer(s.Addr)

	logHandler := olog.NewLogHandler(s.logger, logLine)

	websrv.RouteDefault(http.NotFoundHandler(), orujo.M(logHandler))

	websrv.Route(`^/cmd/list$`,
		http.HandlerFunc(s.listCommandsHandler),
		orujo.M(logHandler))

	if err := websrv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (s *Server) initCommands() {
	s.mutex.Lock()
	s.commands = make(map[string]*command)

	files, err := ioutil.ReadDir(s.CmdDir)
	if err != nil {
		s.logger.Println("command update error:", err)
		return
	}

	for _, f := range files {
		if path.Ext(f.Name()) != ".cmd" {
			continue
		}

		fileName := path.Join(s.CmdDir, f.Name())
		cmd, err := readCommandFile(fileName)
		if err != nil {
			s.logger.Println("command update error:", err)
			return
		}

		s.commands[cmd.Name] = cmd
		s.logger.Println("command updated:", cmd.Name)
	}
	s.mutex.Unlock()
}

func (s *Server) setupWatcher() error {
	if s.CmdDir == "" {
		return errors.New("Server.CmdDir cannot be an empty string")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go s.trackCommands(watcher)

	if err = watcher.Watch(s.CmdDir); err != nil {
		return err
	}

	return nil
}

func (s *Server) trackCommands(watcher *fsnotify.Watcher) {
	for {
		select {
		case ev := <-watcher.Event:
			if path.Ext(ev.Name) != ".cmd" {
				continue
			}
			s.initCommands()
		case err := <-watcher.Error:
			s.logger.Println("watcher error:", err)
		}
	}
}

const logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}
{{range  $err := .Errors}}  Err: {{$err}}
{{end}}`
