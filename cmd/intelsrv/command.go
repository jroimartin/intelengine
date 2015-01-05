// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os/exec"
	"path"
	"strings"
)

const cmdExt = ".cmd"

type command struct {
	Name        string
	Description string
	Cmd         string
	Args        []string
	Class       string
}

func newCommand(filename string) (*command, error) {
	if path.Ext(filename) != cmdExt {
		return nil, errors.New("not a command file")
	}

	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cmd := &command{}
	if err = json.Unmarshal(f, &cmd); err != nil {
		return nil, err
	}
	cmd.Name = strings.TrimSuffix(path.Base(filename), cmdExt)

	return cmd, nil
}

func (cmd *command) exec(r io.Reader) (output []byte, err error) {
	c := exec.Command(cmd.Cmd, cmd.Args...)
	c.Stdin = r
	return c.Output()
}
