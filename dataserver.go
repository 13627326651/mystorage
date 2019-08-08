package main

import (
	"log"
	"net/http"
	"./objects"
	"./hearbeat"
	"./locate"
	"os"
)


func main(){
	http.Handle("/objects/", objects.DataHandler{})
	go locate.StartLocate()
	go hearbeat.StartHearbeat()

	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"), nil))

}