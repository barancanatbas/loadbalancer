package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ReadFile()
	proxy := NewProxy(cfg, "round robine")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	go Checkhealth(exit)

	server := http.Server{
		Addr:    ":" + cfg.Port,
		Handler: http.HandlerFunc(proxy.Handle),
	}
	fmt.Println("listening on port: ", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
