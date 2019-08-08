package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type PutStream struct{
	writer *io.PipeWriter
	c chan error
}

func NewPutStream(server, object string) *PutStream{
	r, w := io.Pipe()
	c := make(chan error)
	go func(){
		request, _ := http.NewRequest("PUT", "http://"+server+"/objects/"+object, r)
		fmt.Println("new request", request.URL)
		client := http.Client{}
		r, e := client.Do(request)
		if e == nil && r.StatusCode != http.StatusOK {
			e = fmt.Errorf("dataserver return status code %d\n", r.StatusCode)
		}
		c <- e
	}()
	return &PutStream{w, c}
}


func (w *PutStream)Write(p []byte)(n int, err error){
	return w.writer.Write(p)
}

func (w *PutStream)Close() error {
	w.writer.Close()
	return <- w.c
}


type GetStream struct{
	reader io.Reader
}


func NewGetStream(server, object string) (*GetStream, error){
	if server == "" || object == "" {
		return nil, fmt.Errorf("invalid dataserver[%s], object[%s]\n ", server, object)
	}

	url := "http://" + server + "/objects/" + object

	return newGetStream(url)
}


func newGetStream(url string) (*GetStream, error){
	r, e := http.Get(url)
	if e != nil {
		return nil, e
	}

	if r.StatusCode != http.StatusOK{
		return nil, fmt.Errorf("dataserver return http code %d\n", r.StatusCode)
	}

	return &GetStream{r.Body}, nil
}

func (r *GetStream)Read(p []byte)(n int, err error){
	return r.reader.Read(p)
}