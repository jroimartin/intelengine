package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/jroimartin/orujo"
	olog "github.com/jroimartin/orujo-handlers/log"
)

var addr = flag.String("addr", ":8080", "HTTP service address")

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "[intelengine] ", log.LstdFlags)
	logHandler := olog.NewLogHandler(logger, logLine)

	s := orujo.NewServer(*addr)

	// TODO: Add routes
	s.RouteDefault(http.NotFoundHandler(), orujo.M(logHandler))

	log.Fatal(s.ListenAndServe())
}

const logLine = `{{.Req.RemoteAddr}} - {{.Req.Method}} {{.Req.RequestURI}}
{{range  $err := .Errors}}  Err: {{$err}}
{{end}}`
