package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ApplicationServer is a simple wrapper around our web service.
// It provides graceful shutdown among other things.
type ApplicationServer struct {
	webserver *http.Server
	Router    *gin.Engine
	basePath  string
}

// NewApplicationServer creates a new instance of our application server.
func NewApplicationServer(
	listenAddr string,
	basePath string,
	logService *LogService,
) ApplicationServer {
	log.Println("NewApplicationServer", "creating server", listenAddr, basePath)

	// configure gin engine
	router := gin.Default()

	// configure routes
	base := normalizePath(basePath)
	router.GET(base+"/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "Hello, World!")
	})
	router.GET(base+"/logs", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"logs": logService.Logs,
		})
	})

	// create application server
	return ApplicationServer{
		basePath: base,
		Router:   router,
		webserver: &http.Server{
			Addr:         listenAddr,
			Handler:      router,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  15 * time.Second,
		},
	}
}

// InstallGETHandler adds new handler with method GET.
func (srv *ApplicationServer) InstallGETHandler(relativePath string, handler gin.HandlerFunc) {
	srv.Router.GET(srv.basePath+normalizePath(relativePath), handler)
}

// InstallPOSTHandler adds new handler with method POST.
func (srv *ApplicationServer) InstallPOSTHandler(relativePath string, handler gin.HandlerFunc) {
	srv.Router.POST(srv.basePath+normalizePath(relativePath), handler)
}

// ListenAndServe starts listening for requests.
func (srv *ApplicationServer) ListenAndServe() error {
	return srv.webserver.ListenAndServe()
}

// GracefulShutdown initiates a graceful shutdown.
func (srv *ApplicationServer) GracefulShutdown(quit <-chan os.Signal, done chan<- bool) {
	<-quit
	log.Println("GracefulShutdown", "server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	srv.webserver.SetKeepAlivesEnabled(false)
	if err := srv.webserver.Shutdown(ctx); err != nil {
		log.Panicln("GracefulShutdown", "could not gracefully shutdown the server", err)
	}

	close(done)
}

// normalizePath ensures that path starts always with a leading slash and has no trailing slash.
// The only exception occurs, when the path is empty. In that case an empty string is returned.
func normalizePath(path string) string {
	if path == "" {
		return path
	}
	// ensure leading slash
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	// remove trailing slash
	path = strings.TrimSuffix(path, "/")
	return path
}
