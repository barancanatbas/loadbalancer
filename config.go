package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var cfg Config

type Config struct {
	Port     string    `json:"port"`
	Backends []Backend `json:"backends"`
	TimeOut  int       `json:"timeOut"`
}

type Backend struct {
	Url      string `json:"url"`
	Health   string `json:"health"`
	IsActive bool   `json:"is_active"`
}

func ReadFile() {
	data, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, &cfg)
	SetActiveAllBackend()
}

func SetActiveAllBackend() {
	for index := range cfg.Backends {
		mu.Lock()
		cfg.Backends[index].IsActive = true
		mu.Unlock()
	}
}

func Checkhealth(exit chan os.Signal) {
	t := time.NewTicker(time.Second * time.Duration(cfg.TimeOut))
	for {
		select {
		case <-exit:
			break
		case <-t.C:
			for index := range cfg.Backends {
				if cfg.Backends[index].Health == "" {
					continue
				}
				status := true
				res, err := http.Get(cfg.Backends[index].Url + cfg.Backends[index].Health)
				if err != nil || res != nil {
					status = false
				} else {
					if res.StatusCode == 200 || res.StatusCode == 201 {
						status = true
					} else {
						status = false
					}
				}

				mu.Lock()
				cfg.Backends[index].IsActive = status
				mu.Unlock()
			}
		}
	}
}
