package main

import (
	"fmt"
	"gocache/gocache"
	"gocache/server"
	"log"
	"net/http"
)

// test data
var db = map[string]string{
	"Tom":  "630",
	"Jack": "589",
	"Sam":  "567",
}

func main() {
	// test data
	gocache.NewGroup("scores", 2<<10, gocache.GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := server.NewHTTPPool(addr)
	log.Println("gocache is runing at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))

}
