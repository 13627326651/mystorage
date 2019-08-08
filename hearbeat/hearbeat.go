package hearbeat

import (
	"../rabbitmq"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time)
var mutex = sync.Mutex{}

func StartHearbeat(){
	r := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer r.Close()
	for{
		r.Publish("apiserver", os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}

func ListenHeartbeat(){
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()

	q.Bind("apiserver")

	go removeExpiredServer()

	c := q.Consume()
	for msg := range c{
		dataServer, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

func removeExpiredServer(){
	for{
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers{
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()

	}
}

func GetDataServers() []string{
	mutex.Lock()
	defer mutex.Unlock()

	servers := make([]string, 0)
	for s, _ := range dataServers{
		servers = append(servers, s)
	}
	return servers
}

func ChooseRandomDataServer() string {
	ds := GetDataServers()
	fmt.Println("get servers", ds)
	n := len(ds)
	if n == 0{
		return ""
	}

	return ds[rand.Intn(n)]
}

