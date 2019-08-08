package locate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)
import "../rabbitmq"

/***************************数据服务**********************/

func Locate(name string) bool {
	_, err := os.Stat(name)

	return !os.IsNotExist(err)

}

func StartLocate(){

	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataserver")
	c := q.Consume()

	for msg := range c{
		object, err := strconv.Unquote(string(msg.Body))
		if err != nil{
			panic(err)
		}

		if Locate(os.Getenv("STORAGE_ROOT") + "/objects/" + object){
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}

/***************************接口服务**********************/
type Handler struct{

}


func (h Handler)ServeHTTP(w http.ResponseWriter, r *http.Request){
	fmt.Printf("recv request:%s\n", r.URL.EscapedPath())
	m := r.Method

	if m != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	o := strings.Split(r.URL.EscapedPath(), "/")[2]
	info := Llocate(o)
	if len(info) == 0{
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, _ := json.Marshal(info)
	w.Write(b)
}

func Llocate(name string) string {
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()

	q.Publish("dataserver", name)

	c := q.Consume()
	go func(){
		time.Sleep(time.Second)

		q.Close()
	}()

	msg := <- c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

func Exist(name string) bool {
	return Llocate(name) != ""
}