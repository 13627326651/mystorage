package main

import (
	"log"
	"net/http"
	"os"
	"./hearbeat"
	"./objects"
	"./locate"
)

func main(){
	go hearbeat.ListenHeartbeat()

	http.Handle("/objects/", objects.InterfaceHandler{})
	http.Handle("/locate/", locate.Handler{})
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))
}
