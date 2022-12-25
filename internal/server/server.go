package server

import (
    "fmt"
    "net/http"
    "time"
)

type Server struct {
    httpServer *http.Server
}

func NewServer() *Server {
    return &Server{}
}

func (s *Server) Run(port uint16, handler http.Handler) error {
    return s.run(port, handler)
}

func (s *Server) run(port uint16, handler http.Handler) error {
    s.httpServer = &http.Server{
        Addr:           fmt.Sprintf(":%d", port),
        Handler:        handler,
        MaxHeaderBytes: 1 << 20, // 1 MB
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
    }

    return s.httpServer.ListenAndServe()
}
