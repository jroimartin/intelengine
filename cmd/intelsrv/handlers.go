// Copyright 2014 The intelengine Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/jroimartin/orujo"
)

const (
	week = 7 * 24 * time.Hour
	sep  = '|'
)

func (s *server) listCommandsHandler(w http.ResponseWriter, r *http.Request) {
	uuid, err := s.client.Call("listCommands", nil, week)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		orujo.RegisterError(w, fmt.Errorf("Call:", err))
		return
	}
	s.logger.Printf("Sent: listCommands(nil) (%v)\n", uuid)
	fmt.Fprintf(w, "{\"UUID\":%q}", uuid)
}

func (s *server) commandResultsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO(jrm): Query DB to get results
	s.mu.Lock()
	defer s.mu.Unlock()
	uuid := strings.TrimPrefix(r.URL.Path, "/cmd/result/")
	result, found := s.results[uuid]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		orujo.RegisterError(w, fmt.Errorf("result not found for uuid: %v", uuid))
		return
	}
	delete(s.results, uuid)
	fmt.Fprint(w, string(result))
}

func (s *server) runCommandHandler(w http.ResponseWriter, r *http.Request) {
	name := strings.TrimPrefix(r.URL.Path, "/cmd/exec/")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		orujo.RegisterError(w, fmt.Errorf("ReadAll:", err))
		return
	}
	data := []byte(fmt.Sprintf("%s%c%s", name, sep, body))
	uuid, err := s.client.Call("execCommand", data, week)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		orujo.RegisterError(w, fmt.Errorf("Call:", err))
		return
	}
	s.logger.Printf("Sent: execCommand(%v) (%v)\n", data, uuid)
	fmt.Fprintf(w, "{\"UUID\":%q}", uuid)
}
