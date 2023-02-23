package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/Dynom/TySug/server"
	"github.com/Dynom/TySug/server/service"
	"github.com/sirupsen/logrus"
)

// Version contains the app version, the value is changed during compile time to the appropriate Git tag
var Version = "dev"

func main() {
	var config Config
	var err error

	config, err = buildConfig("config.toml")
	if err != nil {
		panic(err)
	}

	overrideConfigFromEnv(&config)

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Out = os.Stdout
	logger.Level, err = logrus.ParseLevel(config.Server.Log.Level)

	if err != nil {
		panic(err)
	}

	logger.WithFields(logrus.Fields{
		"version": Version,
		"client":  config.Client,
		"server":  config.Server,
		"CORS":    config.CORS,
	}).Info("Starting up...")

	sr := server.NewServiceRegistry()
	for label, references := range config.References {
		var svc server.Service

		svc, err = service.NewDomain(references, logger)
		if err != nil {
			panic(err)
		}

		sr.Register(label, svc)
	}

	headers := http.Header{}
	for _, h := range config.Server.Headers {
		headers.Add(h.Name, h.Value)
	}

	options := []server.Option{
		server.WithDefaultHeaders(headers),
		server.WithLogger(logger),
		server.WithCORS(config.CORS.AllowedOrigins),
		server.WithInputLimitValidator(config.Client.InputLengthMax),
		server.WithGzipHandler(),
	}

	if config.Server.Profiler.Enable {
		options = append(options, server.WithPProf(config.Server.Profiler.Prefix))
	}

	s := server.NewHTTP(
		sr,
		http.NewServeMux(),
		options...,
	)

	err = s.ListenOnAndServe(config.Server.ListenOn)
	if err != nil {
		panic(err)
	}
}

func buildConfig(fileName string) (Config, error) {
	c := Config{}

	b, err := os.ReadFile(fileName)
	if err != nil {
		return c, fmt.Errorf("unable to open %q, reason: %s", fileName, err)
	}

	_, err = toml.Decode(string(b), &c)
	if err != nil {
		return c, fmt.Errorf("unable to unmarshal %q, reason: %s", fileName, err)
	}

	return c, nil
}

func overrideConfigFromEnv(c *Config) {
	if v, exists := os.LookupEnv("LISTEN_URL"); exists {
		c.Server.ListenOn = v
	}

	if v, exists := os.LookupEnv("LOG_LEVEL"); exists {
		c.Server.Log.Level = v
	}

	if v, exists := os.LookupEnv("PROFILER_PREFIX"); exists {
		c.Server.Profiler.Prefix = v
	}

	if v, exists := os.LookupEnv("PROFILER_ENABLE"); exists {
		c.Server.Profiler.Enable = v == "true"
	}
}
