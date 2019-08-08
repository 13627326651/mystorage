package objects

import (
	"../hearbeat"
	"../objectstream"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"../locate"
)

/********************接口服务*********************/
type InterfaceHandler struct{

}


func (h InterfaceHandler)ServeHTTP(w http.ResponseWriter, r *http.Request){
	m := r.Method
	if m == http.MethodPut {
		iput(w, r)
	}

	if m == http.MethodGet {
		iget(w, r)
	}
}

func iget(w http.ResponseWriter, r *http.Request){
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	stream, e := getStream(object)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
}


func getStream(object string)(*objectstream.GetStream, error){
	server := locate.Llocate(object)
	if server == "" {
		return nil, fmt.Errorf("object[%s] locate fail\n", server)
	}

	return objectstream.NewGetStream(server, object)
}

func iput(w http.ResponseWriter, r *http.Request){

	object := strings.Split(r.URL.EscapedPath(), "/")[2]

	c, e := storeObject(r.Body, object)
	if e != nil{
		log.Println(e)
	}
	w.WriteHeader(c)
}

func storeObject(reader io.Reader, object string) (int, error){

	server := hearbeat.ChooseRandomDataServer()
	if server == ""{
		return http.StatusServiceUnavailable, fmt.Errorf("cannot find any dataserver\n")
	}

	stream := objectstream.NewPutStream(server, object)
	io.Copy(stream, reader)
	e := stream.Close()
	if e != nil{
		return http.StatusInternalServerError, e
	}

	return http.StatusOK, nil
}

/********************数据服务*********************/


type DataHandler struct{

}
func (h DataHandler)ServeHTTP(w http.ResponseWriter, r *http.Request){
	m := r.Method
	if m == http.MethodPut {
		dput(w, r)
	}
	if m == http.MethodGet{
		dget(w, r)
	}
}

func dget(w http.ResponseWriter, r *http.Request){
	f, e := os.Open(os.Getenv("STORAGE_ROOT") + "/objects/" + strings.Split(r.URL.EscapedPath(), "/")[2])
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	io.Copy(w, f)
}

func dput(w http.ResponseWriter, r *http.Request){
	fmt.Printf("dataserver recv put object[%s]\n",strings.Split(r.URL.EscapedPath(), "/")[2])
	f, e := os.Create(os.Getenv("STORAGE_ROOT") + "/objects/" + strings.Split(r.URL.EscapedPath(), "/")[2])
	if e != nil{
		fmt.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	io.Copy(f, r.Body)
}