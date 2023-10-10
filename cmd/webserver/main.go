package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/EricNeid/go-webserver-gin/server"
	"github.com/EricNeid/go-webserver-gin/writer"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logFile     = "logs/webserver.log"
	listenAddr  = ":5000"
	basePath    = ""
	serveStatic = ""
)

func main() {
	// read arguments
	if value, isSet := os.LookupEnv("LISTEN_ADDR"); isSet {
		listenAddr = value
	}
	if value, isSet := os.LookupEnv("BASE_PATH"); isSet {
		basePath = value
	}
	if value, isSet := os.LookupEnv("SERVE_STATIC"); isSet {
		serveStatic = value
	}
	// cli arguments can override environment variables
	flag.StringVar(&listenAddr, "listen-addr", listenAddr, "server listen address")
	flag.StringVar(&basePath, "base-path", basePath, "base path to serve endpoints")
	flag.StringVar(&serveStatic, "serve-static", serveStatic, "serve static files, ie. public=>/dashboard")
	flag.Parse()

	// prepare logging and gin
	logService := server.LogService{
		Max:              5000,
		MaxMessageLength: 250,
	}
	logOut := writer.LazyMultiWriter(
		os.Stdout,
		&logService,
		&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		},
	)
	gin.DefaultWriter = logOut
	log.SetOutput(logOut)
	log.SetPrefix("[APP] ")

	log.Println("main", "starting application")

	// print system proxy
	log.Println("main", "using system proxy")
	log.Println("main", "HTTP_PROXY:", os.Getenv("HTTP_PROXY"))
	log.Println("main", "HTTPS_PROXY:", os.Getenv("HTTPS_PROXY"))

	// prepare graceful shutdown channel
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// create server
	log.Println("main", "creating server")
	gin.SetMode(gin.ReleaseMode)
	srv := server.NewApplicationServer(listenAddr, basePath, &logService)

	if serveStatic != "" {
		root, path, err := parseServeStaticArg(serveStatic)
		if err != nil {
			log.Panicln("main", "could not serve static files", err)
		}
		srv.Router.Static(path, root)
	}

	go srv.GracefulShutdown(quit, done)

	// start listening
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln("main", "could not start listening", err)
	}

	// wait for shutdown
	<-done
	log.Println("main", "server stopped")
}

func parseServeStaticArg(arg string) (root, relativePath string, err error) {
	segments := strings.Split(arg, "=>")
	if len(segments) != 2 {
		return "", "", fmt.Errorf("invalid argument given %s, expecting root=>relativePath", arg)
	}
	return segments[0], segments[1], nil
}
