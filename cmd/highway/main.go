package main

import (
	"fmt"
	"github.com/Highway-Project/highway/config"
	"github.com/Highway-Project/highway/internal/server"
	"github.com/Highway-Project/highway/logging"
	"github.com/creasty/defaults"
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
	logging.Logger.Infof("serving on port :%s", cfg.Global.Port)
	logging.Logger.Fatal(s.ListenAndServe())
}
