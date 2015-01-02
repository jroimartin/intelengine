// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/jroimartin/orujo"
)

func (s *Server) listCommandsHandler(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	b, err := json.Marshal(s.commands)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		orujo.RegisterError(w, fmt.Errorf("cannot marshal commands: %v", err))
		return
	}

	fmt.Fprint(w, string(b))
}

func (s *Server) runCommandHandler(w http.ResponseWriter, r *http.Request) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	name := strings.TrimPrefix(r.URL.Path, "/cmd/exec/")

	cmd := s.command(name)
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
