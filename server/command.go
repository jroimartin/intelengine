// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
)

type command struct {
	Name        string
	Description string
	Cmd         string
	Args        []string
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

// TODO Use stdin as input
func (cmd *command) Run() (output []byte, err error) {
	return exec.Command(cmd.Cmd, cmd.Args...).Output()
}

// TODO return json data
func (s *Server) listCommandsHandler(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	for _, cmd := range s.commands {
		fmt.Fprintf(w, "%#v\n", cmd)
	}
	s.mutex.RUnlock()
}

func (s *Server) getCommand(name string) *command {
	var cmd *command
	s.mutex.RLock()
	for _, c := range s.commands {
		if c.Name == name {
			cmd = c
			break
		}
	}
	s.mutex.RUnlock()
	return cmd
}

// TODO pass body by stdin to command
func (s *Server) runCommandHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/cmd/exec/")

	s.mutex.RLock()
	if cmd := s.getCommand(name); cmd != nil {
		if out, err := cmd.Run(); err == nil {
			fmt.Fprint(w, string(out))
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			s.logger.Println("command execution error:", err)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		s.logger.Println("command not found:", name)
	}
	s.mutex.RUnlock()
}
