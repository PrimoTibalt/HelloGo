package main

import (
	"net/http"
	"os"

	"github.com/fsnotify/fsnotify"
)

func AttachWatcher(path string) error {
	tailscalepartnerip := os.Getenv("TAILSCALE_PARTNER_IP")
	tailscalepartnerurl := tailscalepartnerip + ":8081"
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	watcher.Add(path)
	go func() {
		select {
		case event := <-watcher.Events:
			switch event.Op {
			case fsnotify.Create:
				http.DefaultClient.Post(tailscalepartnerurl, "text", nil)
			case fsnotify.Remove:
			case fsnotify.Rename:
			case fsnotify.Write:
			}
		case err = <-watcher.Errors:
			panic(err)
		}
	}()

	return nil
}
