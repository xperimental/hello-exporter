package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/pflag"
	"github.com/xperimental/hello-exporter/api"
)

type config struct {
	Addr     string
	Username string
	Password string
}

func parseConfig() (config, error) {
	cfg := config{}
	pflag.StringVarP(&cfg.Addr, "addr", "a", ":8080", "Address to listen on.")
	pflag.StringVarP(&cfg.Username, "username", "u", "", "Username of Hello account.")
	pflag.StringVarP(&cfg.Password, "password", "p", "", "Password of Hello account.")
	pflag.Parse()

	if len(cfg.Addr) == 0 {
		return cfg, errors.New("no listen address")
	}

	if len(cfg.Username) == 0 {
		return cfg, errors.New("username can not be blank")
	}

	if len(cfg.Password) == 0 {
		return cfg, errors.New("password can not be blank")
	}

	return cfg, nil
}

func main() {
	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("Error in configuration: %s", err)
	}

	log.Printf("Login as %s", cfg.Username)
	client := api.NewClient(cfg.Username, cfg.Password)

	metrics := newCollector(client)
	prometheus.MustRegister(metrics)

	http.Handle("/metrics", prometheus.UninstrumentedHandler())
	http.Handle("/", http.RedirectHandler("/metrics", http.StatusFound))

	log.Printf("Listen on %s...", cfg.Addr)
	log.Fatal(http.ListenAndServe(cfg.Addr, nil))
}
