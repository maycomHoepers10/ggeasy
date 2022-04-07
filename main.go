// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	//log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")

}

func main() {
	flag.Parse()
	hub := NewHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/wschat", func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r)
	})

	port := os.Getenv("PORT")

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
