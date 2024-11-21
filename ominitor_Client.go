package ominitor

import (
	"embed"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
)

type Client struct {
	opts *redis.Options
}

//go:embed dev/src/*
var embedSource embed.FS

func (c *Client) Start(address string) {
	c.listen(address, true)
}

func (c *Client) Develop(address string) {
	c.listen(address, false)
}

func (c *Client) listen(address string, embedModel bool) {

	manager := NewManager(c.opts)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		part := strings.Split(r.URL.Path, "/")
		if r.URL.Path != "/" && len(part) > 1 && len(strings.Split(part[1], ".")) == 1 {
			manager.Handler(w, r)
			return
		}

		filePath := r.URL.Path
		if r.URL.Path == "/" {
			filePath = "/index.html"
		}
		var data []byte
		if embedModel {
			filePath = "dev/src" + filePath
			data, _ = embedSource.ReadFile(filePath)
		} else {
			filePath = "src" + filePath
			data, _ = os.ReadFile(filePath)
		}
		w.Write(data)
	})

	log.Println("omi web manager server is running on http://" + address)

	http.ListenAndServe(":"+strings.Split(address, ":")[1], nil)
}
