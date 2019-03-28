package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/adamb/netpupper/controller"
	"github.com/adamb/netpupper/errors"
	"io/ioutil"
	"log"
	"net/http"
)

type API struct {
	Controller *controller.Controller
}

type APIClient struct {
}

func StartServerApi() API {
	con := controller.Controller{}
	a := API{
		Controller: &con,
	}

	http.HandleFunc("/register", a.register)
	log.Fatal(http.ListenAndServe(":8999", nil))

	return a
}

/*
Receive a register from a client
*/
func (a *API) register(w http.ResponseWriter, r *http.Request) {
	reg := controller.Client{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &reg)

	a.Controller.AddClient(reg)

	m := Message{
		Value: "OK",
	}
	b, err := json.Marshal(m)
	errors.CheckError(err)
	w.Write(b)

}

func (a *APIClient) SendRegister() {
	gt := controller.Tag{
		Name:  "Geoloc",
		Value: "1",
	}
	tags := []controller.Tag{gt}
	reg := controller.Client{
		Addr: "1.1.1.1",
		Tags: tags,
	}
	b, err := json.Marshal(reg)
	errors.CheckError(err)
	resp, err := http.Post("http://127.0.0.1:8999/register", "application/json", bytes.NewBuffer(b))
	errors.CheckError(err)

	msg := Message{}
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &msg)
	fmt.Printf("DEBUG: %v\n", msg.Value)
}
