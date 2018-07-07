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

var config = Config{}

func main() {
	var err error

	config, err = buildConfig("config.yml")

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
		"client": config.Client,
		"server": config.Server,
		"CORS":   config.CORS,
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
		fmt.Printf("Unable to open 'config.yml', reason: %s\n%s", err, b)
	}

	err = yaml.Unmarshal(b, &c)
	if err != nil {
		fmt.Printf("Unable to unmarshal 'config.yml', reason: %s\n%s", err, b)
	}

	return c, nil
}
