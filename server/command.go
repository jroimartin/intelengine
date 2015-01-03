// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os/exec"
)

type command struct {
	Name        string
	Description string
	Cmd         string
	Args        []string
	Class       string
}

func newCommand(filename string) (*command, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cmd := &command{}
	if err = json.Unmarshal(f, &cmd); err != nil {
		return nil, err
	}

	if cmd.Name == "" {
		return nil, errors.New("Command name cannot be an empty string")
	}

	return cmd, nil
}

func (cmd *command) exec(r io.Reader) (output []byte, err error) {
	c := exec.Command(cmd.Cmd, cmd.Args...)
	c.Stdin = r
	return c.Output()
}
