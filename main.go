package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println("en error occured:", err.Error())
			return
		}
		fmt.Println(string(b))
	})
	server := &http.Server{
		Addr:        "0.0.0.0:3455",
		Handler:     mux,
		ReadTimeout: time.Second * 2,
	}
	server.ListenAndServe()
}
