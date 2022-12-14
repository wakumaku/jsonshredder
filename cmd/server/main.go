package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wakumaku/jsonshredder/cmd/server/handler"
	"wakumaku/jsonshredder/internal/config"
	"wakumaku/jsonshredder/internal/service"

	"github.com/rs/zerolog"
)

// Release info, overrided by ldflags
var (
	Date    string = "today-dev"
	Version string = "0.0.0-dev"
	Commit  string = "00FF0-dev"
)

func main() {
	// Parse input params
	var (
		p        string
		showHelp bool
	)
	flag.StringVar(&p, "config", "", "path to config file")
	flag.BoolVar(&showHelp, "help", false, "shows this help")
	flag.Parse()

	if showHelp {
		fmt.Printf(`Version: %s\nDate: %s\nCommit:%s\n`,
			Version, Date, Commit)
		flag.Usage()
		os.Exit(0)
	}

	if p == "" {
		flag.Usage()
		os.Exit(1)
	}

	// Loads config
	cfg, err := config.LoadFromFile(p)
	if err != nil {
		fmt.Println("ERR: ", err)
		os.Exit(1)
	}

	// initializes context
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)
	defer cancel()

	zerolog.SetGlobalLevel(cfg.LogLevel)
	logger := zerolog.New(os.Stdout).With().
		Str("appname", "jsonshredder").
		Str("version", Version).
		Timestamp().
		Logger()
	logger.Debug().Msg("Configuration loaded")

	// Builds Shredder service
	shredSrv := service.NewShredder(cfg.Transformations, &logger)
	ffwSrv := service.NewForwarder(ctx,cfg.Forwarders, &logger)

	// Initializes the HTTP server
	connTimeout := 5 * time.Second
	server := New(cfg.Port, handler.Router(shredSrv, ffwSrv, &logger), connTimeout)

	go func() {
		<-ctx.Done()
		logger.Warn().Msg("closing server!")
		cancel()
	}()

	logger.Info().Str("section", "main").Msg("Starting server ...")
	logger.Fatal().Str("section", "main").Err(server.Run(ctx)).Send()
	time.Sleep(time.Second)
}
