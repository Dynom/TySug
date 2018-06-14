package main

import (
	"net/http"
	"os"

	"github.com/Dynom/TySug/server"
	"github.com/Dynom/TySug/server/service"
	"github.com/sirupsen/logrus"

	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config holds TySug's central config parameters
type Config struct {
	References map[string][]string `yaml:"references"`
	Client     struct {
		ReferencesMax  int `yaml:"referencesMax"`
		InputLengthMax int `yaml:"inputLengthMax"`
	}
	CORS struct {
		AllowedOrigins []string `yaml:"allowedOrigins"`
	}
	Server struct {
		ListenOn string `yaml:"listenOn"`
	}
}

var config = Config{}

func main() {

	config, err := buildConfig("config.yml")

	if err != nil {
		panic(err)
	}

	logger := logrus.New()
	logger.Formatter = &logrus.JSONFormatter{}
	logger.Info("Starting up...")
	logger.Level = logrus.DebugLevel
	logger.Out = os.Stdout

	sr := server.NewServiceRegistry()
	for label, references := range config.References {
		svc, err := service.NewDomain(references, logger)
		if err != nil {
			panic(err)
		}

		sr.Register(label, svc)
	}

	mux := http.NewServeMux()
	s := server.NewHTTP(sr, mux, server.WithLogger(logger))

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
