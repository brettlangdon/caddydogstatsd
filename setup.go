package caddydogstatsd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/datadog/datadog-go/statsd"
	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func init() {
	// register our plugin with Caddy
	caddy.RegisterPlugin("dogstatsd", caddy.Plugin{
		ServerType: "http",
		Action:     setup,
	})
}

func setup(c *caddy.Controller) error {
	for c.Next() {
		// only parse if the initial directive is "dogstatsd"
		if c.Val() != "dogstatsd" {
			continue
		}

		// default config values
		var namespace = ""
		var host = "127.0.0.1:8125"
		var globalTags = []string{}
		var sampleRate = 1.0

		// if we have a block, then parse that
		// e.g.
		//   dogstatsd {
		//     host 127.0.0.1:8125
		//   }
		for c.NextBlock() {
			// each line if of the format `{key} {arg} [{arg}...]`
			var key string
			key = c.Val()

			var args []string
			args = c.RemainingArgs()
			// we expect every directive to have at least 1 argument
			if len(args) == 0 {
				return c.ArgErr()
			}

			// parse directives
			switch key {
			case "host":
				host = args[0]
			case "samplerate":
				var err error
				sampleRate, err = strconv.ParseFloat(args[0], 64)
				if err != nil {
					return c.SyntaxErr(fmt.Sprintf("expected float for \"samplerate\", instead found \"%s\"", args[0]))
				}
			case "namespace":
				namespace = args[0]
				if !strings.HasSuffix(namespace, ".") {
					namespace += "."
				}
			case "tags":
				globalTags = args
			default:
				return c.SyntaxErr(fmt.Sprintf("expected one of \"host\", \"samplerate\", \"namespace\", \"tags\", instead found \"%s\"", key))
			}
		}

		// handle non-block configuration
		// e.g.
		//   dogstatsd [{host:port} [{samplerate}]]
		if c.NextArg() {
			var args []string
			args = c.RemainingArgs()
			if len(args) > 0 {
				host = args[0]
			}
			if len(args) > 1 {
				var err error
				sampleRate, err = strconv.ParseFloat(args[1], 64)
				if err != nil {
					return c.SyntaxErr(fmt.Sprintf("expected float for \"samplerate\", instead found \"%s\"", args[1]))
				}
			}
		}

		// add our middleware
		var cfg *httpserver.SiteConfig
		cfg = httpserver.GetConfig(c)
		cfg.AddMiddleware(func(next httpserver.Handler) httpserver.Handler {
			var client *statsd.Client
			var err error
			client, err = statsd.New(host)
			if err == nil {
				client.Namespace = namespace
				client.Tags = globalTags
			}

			return dogstatsdHandler{
				Client:     client,
				SampleRate: sampleRate,
				Next:       next,
			}
		})
	}
	return nil
}
