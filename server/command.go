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
)

type command struct {
	Name        string
	Description string
	Path        string
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

func (cmd *command) Run() (output string, err error) {
	return "", nil
}

// TODO return json data
func (s *Server) listCommandsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, s.commands)
}
