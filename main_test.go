package main

import (
	"testing"

	"github.com/BurntSushi/toml"
)

const tomlTestStr = `
[client]
  # The maximum length (in bytes) of the "input" field, to accept to the finder.
  inputLengthMax = 50

[CORS]
  # Not defining allowedOrigins will use wildcard '*', meaning all clients are allowed.
  allowedOrigins = [
    "http://example.org"
  ]

[server]
  # The interface and port to listen on, "0.0.0.0" means all interfaces
  listenOn = "0.0.0.0:1337"
  [server.log]
    # The minimum logging level
    level = "info"
  [server.profiler]
    enable = true
    prefix = "debug"


# The list of references
[references]
  test = [
    "beek",
    "beel",
    "been",
    "bear",
    "bare",
    "beer",
    "bool",
    "boot",
  ]
`

func TestBuildConfig(t *testing.T) {
	var c Config

	_, err := toml.Decode(tomlTestStr, &c)
	if err != nil {
		t.Error(err)
	}

	if c.Client.InputLengthMax != 50 {
		t.Errorf("Expected InputLengthMax to be 50, instead it was %d", c.Client.InputLengthMax)
	}

	if c.Server.Profiler.Enable == false {
		t.Error("Expected the profile to be enabled")
	}

	if l := len(c.References["test"]); l != 8 {
		t.Errorf("Expected 8 items in the test reference list, instead it was%d", l)
	}
}
