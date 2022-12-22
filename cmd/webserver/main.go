package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/EricNeid/go-webserver-gin/server"
	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	logFile    string = "logs/webserver.log"
	listenAddr string = ":5000"
	basePath   string = ""
)

func main() {
	// read arguments
	if value, isSet := os.LookupEnv("LISTEN_ADDR"); isSet {
		listenAddr = value
	}
	if value, isSet := os.LookupEnv("BASE_PATH"); isSet {
		basePath = value
	}
	// cli arguments can override environment variables
	flag.StringVar(&listenAddr, "listen-addr", listenAddr, "server listen address")
	flag.StringVar(&basePath, "base-path", basePath, "base path to serve endpoints")
	flag.Parse()

	// print system proxy
	log.Println("main", "using system proxy")
	log.Println("main", "HTTP_PROXY:", os.Getenv("HTTP_PROXY"))
	log.Println("main", "HTTPS_PROXY:", os.Getenv("HTTPS_PROXY"))

	// prepare logging and gin
	logService := server.LogService{
		Max: 5000,
	}
	logOut := LazyMultiWriter(
		os.Stdout,
		&logService,
		&lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    500, // megabytes
			MaxBackups: 3,
			MaxAge:     28, //days
		},
	)
	gin.DefaultWriter = logOut
	log.SetOutput(logOut)
	log.SetPrefix("[APP] ")

	// prepare gracefull shutdown channel
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	// create server
	log.Println("main", "creating server")
	gin.SetMode(gin.ReleaseMode)
	srv := server.NewApplicationServer(listenAddr, basePath, &logService)

	go srv.GracefullShutdown(quit, done)

	// start listening
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalln("main", "could not start listening", err)
	}

	// wait for shutdown
	<-done
	log.Println("main", "server stopped")
}
