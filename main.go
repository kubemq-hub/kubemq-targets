package main

import (
	"context"
	"github.com/kubemq-hub/kubemq-target-connectors/binding"
	"github.com/kubemq-hub/kubemq-target-connectors/pkg/logger"

	"github.com/kubemq-hub/kubemq-target-connectors/config"

	"os"
	"os/signal"
	"syscall"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)
var (
	bindingMap = map[string]*binding.Binder{}
)
var (
	log *logger.Logger
)

func start(ctx context.Context, cfg *config.Config) error {
	for _, bindingCfg := range cfg.Bindings {
		binder := binding.New()
		err := binder.Init(ctx, bindingCfg)
		if err != nil {
			return err
		}
		err = binder.Start(ctx)
		if err != nil {
			return err
		}
		bindingMap[bindingCfg.Name] = binder
	}
	return nil
}

func stop() {
	for _, binder := range bindingMap {
		_ = binder.Stop()
	}
}

func run() error {
	var gracefulShutdown = make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	err = cfg.Validate()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = start(ctx, cfg)
	if err != nil {
		return err
	}
	<-gracefulShutdown
	stop()
	return nil
}
func main() {
	log = logger.NewLogger("main")
	log.Infof("starting kubemq targets connector version: %s, commit: %s, date %s", version, commit, date)
	if err := run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
