package main

import (
	"net/http"

	"github.com/Dynom/TySug/server"
	"github.com/Dynom/TySug/server/service"
	"github.com/sirupsen/logrus"

	"fmt"
	"io/ioutil"

	"os"

	"gopkg.in/yaml.v2"
)

// Config holds TySug's central config parameters
type Config struct {
	References map[string][]string `yaml:"references"`
	Client     struct {
		//ReferencesMax  int `yaml:"referencesMax"`
		InputLengthMax int `yaml:"inputLengthMax"`
	} `yaml:"client"`
	CORS struct {
		AllowedOrigins []string `yaml:"allowedOrigins"`
	} `yaml:"CORS"`
	Server struct {
		ListenOn string `yaml:"listenOn"`
		Log      struct {
			Level string `yaml:"level"`
		} `yaml:"log"`
		Profiler struct {
			Enable bool   `yaml:"enable"`
			Prefix string `yaml:"prefix"`
		} `yaml:"profiler"`
	} `yaml:"server"`
}

// Version contains the app version, the value is changed during compile time to the appropriate Git tag
var Version = "dev"

func main() {
	var config Config
	var err error

	config, err = buildConfig("config.yml")

	if err != nil {
		panic(err)
	}

	err = overrideConfigFromEnv(&config)
	if err != nil {
		panic(err)
	}

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

	options := []server.Option{
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

	b, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Unable to open '%s', reason: %s\n%s", fileName, err, b)
	}

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		fmt.Printf("Unable to unmarshal '%s', reason: %s\n%s", fileName, err, b)
	}

	return c, nil
}

func overrideConfigFromEnv(c *Config) error {
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
		if v == "true" {
			c.Server.Profiler.Enable = true
		} else {
			c.Server.Profiler.Enable = false
		}
	}

	return nil
}
