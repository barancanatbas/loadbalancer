package main

import "net/http"

type IProxy interface {
	Handle(w http.ResponseWriter, r *http.Request)
}

func NewProxy(cfg Config, strategy string) IProxy {
	if strategy == "round robine" {
		return &RoundRobine{Config: cfg}
	}
	// default olarak round robine strategy
	return &RoundRobine{Config: cfg}
}
