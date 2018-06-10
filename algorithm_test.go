package main_test

import (
	"testing"

	"github.com/Dynom/TySug/finder"
	"github.com/xrash/smetrics"
	"github.com/zikes/sift4"
)

// Several algorithms to test with.
var algorithms = map[string]finder.AlgWrapper{
	"Ukkonen 1/1/1": func(a, b string) float64 {
		return -1 * float64(smetrics.Ukkonen(a, b, 1, 1, 1))
	},
	"JaroWinkler .7/4": func(a, b string) float64 {
		return smetrics.JaroWinkler(a, b, .7, 4)
	},
	"WagnerFischer 1/1/1": func(a, b string) float64 {
		return -1 * float64(smetrics.WagnerFischer(a, b, 1, 1, 1))
	},
	"Sift4": func(a, b string) float64 {
		return -1 * sift4.New().Distance(a, b)
	},
}

// The order is important. Put the most frequently used reference first.
var domains = []string{
	"gmail.com", "yahoo.com", "hotmail.com", "aol.com", "hotmail.co.uk", "hotmail.fr", "msn.com", "yahoo.fr",
	"wanadoo.fr", "orange.fr", "comcast.net", "yahoo.co.uk", "yahoo.com.br", "yahoo.co.in", "live.com",
	"rediffmail.com", "free.fr", "gmx.de", "web.de", "yandex.ru", "ymail.com", "libero.it", "outlook.com",
	"uol.com.br", "bol.com.br", "mail.ru", "cox.net", "hotmail.it", "sbcglobal.net", "sfr.fr", "live.fr",
	"verizon.net", "live.co.uk", "googlemail.com", "yahoo.es", "ig.com.br", "live.nl", "bigpond.com",
	"terra.com.br", "yahoo.it", "neuf.fr", "yahoo.de", "alice.it", "rocketmail.com", "att.net", "laposte.net",
	"facebook.com", "bellsouth.net", "yahoo.in", "hotmail.es", "charter.net", "yahoo.ca", "yahoo.com.au",
	"rambler.ru", "hotmail.de", "tiscali.it", "shaw.ca", "yahoo.co.jp", "sky.com", "earthlink.net", "optonline.net",
	"freenet.de", "t-online.de", "aliceadsl.fr", "virgilio.it", "home.nl", "qq.com", "telenet.be", "me.com",
	"yahoo.com.ar", "tiscali.co.uk", "yahoo.com.mx", "voila.fr", "gmx.net", "mail.com", "planet.nl", "tin.it",
	"live.it", "ntlworld.com", "arcor.de", "yahoo.co.id", "frontiernet.net", "hetnet.nl", "live.com.au",
	"yahoo.com.sg", "zonnet.nl", "club-internet.fr", "juno.com", "optusnet.com.au", "blueyonder.co.uk",
	"bluewin.ch", "skynet.be", "sympatico.ca", "windstream.net", "mac.com", "centurytel.net", "chello.nl",
	"live.ca", "aim.com", "bigpond.net.au",
	"hotmail.nl", "ziggo.nl", "live.com",
}

// TestAlgorithms generates a result for the algorithms defined above. It never errors out, use -test.v to have it show
// up in the go test output
func TestAlgorithms(t *testing.T) {
	testData := map[string][]string{
		// Expected - []spelling mistakes
		"hotmail.com": {"hotmail.co", "homail.com", "hotmal.com", "hotmai.com", "hotmailcom", "hotmal.co", "hotmai.com/", "hotmaol.com", "hotmail.con", "hormail.com", "hotnail.com", "hotmaul.com"},
		"hotmail.nl":  {"hotmail.bl", "hotmal.nl", "hotmai.nl", "hotmailnl", "hotmal.nl", "hotmai.nl/"},
		"gmail.com":   {"gmai.com", "gmail.dom", "gnail.com", "gamil.com", "hmail.com", "gmail.con", "gmail.co", "email.com", "hmail.com"},
		"mail.com":    {"maill.com", "mail.co", "mail.com/", "mail..com"},
		"live.com":    {"life.com"},
		"ziggo.nl":    {"zigo.nl"},

		// Failures
		/*
			"gmail.com": {
				"gail.com", // gail.com now leads to mail.com while I think there is a higher possibility that people
							// mean gmail.com on an English keyboard

			},
		*/
	}

	for name, alg := range algorithms {
		t.Run(name, func(t *testing.T) {
			sug, _ := finder.New(domains, finder.OptSetAlgorithm(alg))

			// Running combination tests for each domain, against our reference list.
			for expectedDomain, emailsToTest := range testData {
				for _, domain := range emailsToTest {

					bestMatch, score := sug.Find(domain)
					if bestMatch != expectedDomain {
						t.Logf("Related score: %f", score)
						t.Logf("Expected '%s' to result in '%s'. Instead I got: '%s'.", domain, expectedDomain, bestMatch)
					}
				}
			}
		})
	}
}
