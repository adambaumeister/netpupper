package api

import (
	"encoding/json"
	"github.com/adamb/netpupper/controller"
	"io/ioutil"
	"log"
	"net/http"
)

type API struct {
}

func StartServerApi() API {
	a := API{}
	http.HandleFunc("/register", register)
	log.Fatal(http.ListenAndServe(":8999", nil))

	return a
}

/*
Register a new client.
*/
func register(w http.ResponseWriter, r *http.Request) {
	reg := controller.Client{}
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &reg)

}
