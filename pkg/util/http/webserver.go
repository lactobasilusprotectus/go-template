package http

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/lactobasilusprotectus/go-template/pkg/common/config"
	"log"
	"net/http"
	"time"
)

// Server represents http server
type Server struct {
	port   string
	server http.Server

	gin *gin.Engine
}

// RouterHandler wraps API handler interface
type RouterHandler interface {
	Register(*gin.Engine)
}

// NewServer constructs new HttpServer
func NewServer(opt config.HttpConfig) *Server {
	//Set gin mode
	if opt.Env == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Creates a router without any middleware by default
	g := gin.New()

	//==================================================================================================
	// Global Middleware
	//==================================================================================================

	// Cors middleware will handle the OPTIONS method and add the corresponding headers to the response
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	g.Use(cors.Default(), gin.Logger(), gin.Recovery())

	return &Server{
		port: opt.Port,
		server: http.Server{
			Addr:        fmt.Sprintf("0.0.0.0:%s", opt.Port),
			IdleTimeout: time.Second * 60,
			ReadTimeout: time.Second * 60,
			Handler:     g,
		},
		gin: g,
	}
}

// RegisterHandler registers our API handler
func (s *Server) RegisterHandler(api RouterHandler) {
	api.Register(s.gin)
}

// Run runs Gin http server on the given port
func (s *Server) Run(env string) {
	go func() {
		log.Println("starting", env, "server on port", s.port)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("err: %s\n", err)
		}
	}()
}

// Stop Shutdown shuts down Gin http server
func (s *Server) Stop() {
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	if err := s.server.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
}
