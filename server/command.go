// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/jroimartin/orujo"
)

type command struct {
	Name        string
	Description string
	Cmd         string
	Args        []string
	Class       string
}

func readCommandFile(path string) (*command, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	cmd := &command{}
	if err = json.Unmarshal(f, &cmd); err != nil {
		return nil, err
	}

	if cmd.Name == "" {
		errors.New("Command name cannot be an empty string")
	}

	return cmd, nil
}

func (cmd *command) exec(r io.Reader) (output []byte, err error) {
	c := exec.Command(cmd.Cmd, cmd.Args...)
	c.Stdin = r
	return c.Output()
}

func (s *Server) listCommandsHandler(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	b, err := json.Marshal(s.commands)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		orujo.RegisterError(w, fmt.Errorf("Cannot marshal commands: %v", err))
		return
	}

	fmt.Fprint(w, string(b))
}

func (s *Server) getCommand(name string) *command {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for _, cmd := range s.commands {
		if cmd.Name == name {
			return cmd
		}
	}
	return nil
}

func (s *Server) runCommandHandler(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	name := strings.TrimPrefix(r.URL.Path, "/cmd/exec/")

	cmd := s.getCommand(name)
	if cmd == nil {
		w.WriteHeader(http.StatusInternalServerError)
		orujo.RegisterError(w, fmt.Errorf("command not found: %v", name))
		return
	}

	out, err := cmd.exec(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		orujo.RegisterError(w, fmt.Errorf("command execution error: %v", err))
		return
	}

	fmt.Fprint(w, string(out))
}
