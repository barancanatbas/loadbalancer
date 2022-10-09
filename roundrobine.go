package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var count = 0
var mu sync.Mutex

type RoundRobine struct {
	Config Config `json:"config"`
}

func (p *RoundRobine) Handle(w http.ResponseWriter, r *http.Request) {
	backendLen := len(p.Config.Backends)
	mu.Lock()
	var goalBackend Backend
	for {
		goalBackend = p.Config.Backends[count%backendLen]
		fmt.Println(goalBackend)
		if !goalBackend.IsActive {
			count++
		} else {
			break
		}
	}
	count++
	goalUrl, err := url.Parse(goalBackend.Url)
	if err != nil {
		fmt.Println(err)
	}

	mu.Unlock()
	reverseProxy := httputil.NewSingleHostReverseProxy(goalUrl)
	reverseProxy.ServeHTTP(w, r)
}
