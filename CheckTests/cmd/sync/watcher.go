package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type CreateTopic struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func AttachWatcher(path string, breaker chan bool) error {
	tailscalepartnerip := os.Getenv("TAILSCALE_PARTNER_IP")
	tailscalepartnerurl := "http://" + tailscalepartnerip + ":8081"
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	watcher.Add(path)
	go func(breaker chan bool) {
		defer watcher.Close()
		for {
			renamed := false
			select {
			case <-breaker:
				fmt.Println("Waiting for operation to complete")
				goto anchor
			case event := <-watcher.Events:
				switch event.Op {
				case fsnotify.Create:
					createTopicModel := CreateTopic{}
					fmt.Println("Creating...")
					createTopicModel.Name = event.Name
					if renamed {
						fmt.Println("And filling...")
						f, err := os.ReadFile(event.Name)
						if err != nil {
							fmt.Println(err)
							continue
						}
						createTopicModel.Content = string(f)
					}
					body, err := json.Marshal(&createTopicModel)
					if err != nil {
						fmt.Printf("could not marshal %v\n", createTopicModel)
					}
					_, err = http.DefaultClient.Post(
						tailscalepartnerurl+"/createTopic",
						"application/json",
						bytes.NewReader(body),
					)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Finished creating")
				case fsnotify.Remove:
					fmt.Println("Removing...")
					_, err := http.DefaultClient.Post(
						tailscalepartnerurl+"/removeTopic",
						"text",
						strings.NewReader(event.Name),
					)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Finished removing")
				case fsnotify.Rename:
					fmt.Println("Renaming...")
					_, err := http.DefaultClient.Post(
						tailscalepartnerurl+"/removeTopic",
						"text",
						strings.NewReader(event.Name),
					)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Finished renaming")
				case fsnotify.Write:
					fmt.Println("Writing...")
					createTopicModel := CreateTopic{}
					createTopicModel.Name = event.Name
					f, err := os.ReadFile(event.Name)
					if err != nil {
						fmt.Println("could not read file " + event.Name)
						continue
					}
					createTopicModel.Content = string(f)
					body, err := json.Marshal(&createTopicModel)
					if err != nil {
						fmt.Println("could not marshal new changes")
						continue
					}
					_, err = http.DefaultClient.Post(
						tailscalepartnerurl+"/changeTopic",
						"application/json",
						bytes.NewReader(body),
					)
					if err != nil {
						fmt.Println(err)
					}
					fmt.Println("Finished writing")
				}
			case err = <-watcher.Errors:
				panic(err)
			}
		}
	anchor:
	}(breaker)

	return nil
}
