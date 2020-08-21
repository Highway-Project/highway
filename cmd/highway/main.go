package main

import (
	"context"
	"fmt"
	"github.com/Highway-Project/highway/config"
	"github.com/Highway-Project/highway/internal/server"
	"github.com/Highway-Project/highway/logging"
	"github.com/creasty/defaults"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println(`
  _    _ _       _                        
 | |  | (_)     | |                       
 | |__| |_  __ _| |____      ____ _ _   _ 
 |  __  | |/ _' | '_ \ \ /\ / / _' | | | |
 | |  | | | (_| | | | \ V  V / (_| | |_| |
 |_|  |_|_|\__, |_| |_|\_/\_/ \__,_|\__, |
            __/ |                    __/ |
           |___/                    |___/`)
	logging.InitLogger("debug", true)

	cfg, err := config.ReadConfig()
	if err != nil {
		logging.Logger.WithError(err).Fatal("could not load config")
	}

	err = defaults.Set(cfg)
	if err != nil {
		logging.Logger.WithError(err).Fatal("could not set default values")
	}

	err = cfg.Validate()
	if err != nil {
		logging.Logger.WithError(err).Fatal("invalid config")
	}

	s, err := server.NewServer(cfg.Global, cfg.RouterSpec, cfg.ServicesSpecs, cfg.RulesSpecs, cfg.MiddlewaresSpecs)
	if err != nil {
		logging.Logger.WithError(err).Fatal("could not create server")
	}

	go func() {
		logging.Logger.Infof("started serving on port :%d", cfg.Global.Port)
		logging.Logger.Fatal(s.ListenAndServe())
	}()

	shutdown := make(chan os.Signal)

	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	<-shutdown

	logging.Logger.Info("shutting down highway gracefully")

	s.SetKeepAlivesEnabled(false)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = s.Shutdown(ctx)
	if err != nil {
		logging.Logger.WithError(err).Error("could not shutdown highway gracefully")
	}

	logging.Logger.Info("exiting highway")

}
