package main

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/mailgun/groupcache/v2"
)

//go:embed templates/*
var files embed.FS

type User struct {
	ID       string
	User     string
	Instance string
}

var group = groupcache.NewGroup("users", 3<<20, groupcache.GetterFunc(
	func(_ context.Context, id string, dest groupcache.Sink) error {
		me, err := os.Hostname()
		if err != nil {
			panic("Get Hostname: " + err.Error())
		}

		log.Printf("Create user-%s on instance %s", id, me)

		user := User{
			ID:       id,
			User:     fmt.Sprintf("user-%s", id),
			Instance: me,
		}

		bs, err := json.Marshal(user)
		if err != nil {
			log.Fatal("Marshal: " + err.Error())
		}

		// Set the user in the groupcache to expire after one minute
		return dest.SetBytes(bs, time.Now().Add(time.Minute))
	},
))

func indexHandler() func(http.ResponseWriter, *http.Request) {
	bs, err := files.ReadFile("templates/index.html")
	if err != nil {
		panic("Reading template: " + err.Error())
	}

	index := template.Must(template.New("index").Parse(string(bs)))

	return func(w http.ResponseWriter, r *http.Request) {
		var user User

		if err := r.ParseForm(); err == nil {
			user.ID = r.PostForm.Get("id")

			if user.ID != "" {
				var bs []byte

				if err := group.Get(r.Context(), user.ID, groupcache.AllocatingByteSliceSink(&bs)); err != nil {
					log.Fatal("Groupcache get: " + err.Error())
				}

				if err := json.Unmarshal(bs, &user); err != nil {
					log.Fatal("Unmarshal: " + err.Error())
				}
			}
		}

		if err := index.Execute(w, user); err != nil {
			log.Fatal("Render: " + err.Error())
		}
	}
}

func newPool(peers []string) *groupcache.HTTPPool {
	pool := groupcache.NewHTTPPoolOpts(peers[0], nil)
	pool.Set(peers...)

	return pool
}

func getPeers() []string {
	me, err := os.Hostname()
	if err != nil {
		panic("Get Hostname: " + err.Error())
	}

	me = fmt.Sprintf("http://%s:8080", me)

	peers := []string{
		"http://app1:8080",
		"http://app2:8080",
		"http://app3:8080",
	}

	for i, v := range peers {
		if v == me {
			peers = append(peers[:i], peers[i+1:]...)
		}
	}

	return append([]string{me}, peers...)
}

func main() {
	groupcache.SetLoggerFromLogger(newLogger())

	peers := getPeers()

	log.Printf("listening on %v", peers[0])
	log.Printf("peers: %v", peers)

	http.HandleFunc("/", indexHandler())
	http.Handle("/_groupcache/", newPool(peers))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
