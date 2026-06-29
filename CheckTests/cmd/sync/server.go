package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
)

func RunServer() {
	mytailscaleip := os.Getenv("TAILSCALE_IP")
	if mytailscaleip == "" {
		panic(errors.New("no ip in TAILSCALE_IP was provided"))
	}
	mux := http.NewServeMux()
	mux.HandleFunc("POST /createTopic", createNewTopic)
	http.ListenAndServe(mytailscaleip+":8081", mux)
}

func createNewTopic(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello")
}
