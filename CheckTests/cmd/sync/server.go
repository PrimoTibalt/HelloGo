package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	mux.HandleFunc("POST /removeTopic", removeTopic)
	mux.HandleFunc("POST /changeTopic", changeTopic)
	http.ListenAndServe(mytailscaleip+":8081", mux)
}

func changeTopic(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	createTopicModel, err := getcreatetopicmodel(r, rw)
	if err != nil {
		return
	}
	f, err := os.OpenFile(createTopicModel.Name, os.O_WRONLY|os.O_CREATE, 0o776)
	if err != nil {
		fmt.Println("could not open file " + createTopicModel.Name)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	_, err = f.Write([]byte(createTopicModel.Content))
	if err != nil {
		fmt.Println("could not write to file " + createTopicModel.Name)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusOK)
}

func getcreatetopicmodel(r *http.Request, rw http.ResponseWriter) (CreateTopic, error) {
	defer r.Body.Close()
	createTopicModel := CreateTopic{}
	allBytes, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		fmt.Println("could not create new topic")
		rw.WriteHeader(http.StatusInternalServerError)
		return createTopicModel, errors.New("Heh")
	}
	err = json.Unmarshal(allBytes, &createTopicModel)
	if err != nil {
		fmt.Println(err)
		fmt.Println("could not unmarchal payload")
		rw.WriteHeader(http.StatusInternalServerError)
		return createTopicModel, errors.New("Heh")
	}
	return createTopicModel, nil
}

func createNewTopic(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	createTopicModel, err := getcreatetopicmodel(r, rw)
	if err != nil {
		return
	}
	f, err := os.Create(createTopicModel.Name)
	if err != nil {
		fmt.Println("could not create topic")
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	defer f.Close()
	_, err = f.Write([]byte(createTopicModel.Content))
	if err != nil {
		fmt.Println("could not write into created file")
		rw.WriteHeader(http.StatusInternalServerError)
	}
	rw.WriteHeader(http.StatusOK)
}

func removeTopic(rw http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println("could not read body")
		fmt.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
	}
	filename := string(body)
	err = os.Remove(filename)
	if err != nil {
		fmt.Println(filename + " does not exist in the system")
	}
	rw.WriteHeader(http.StatusOK)
}
