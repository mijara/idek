package idek

import (
	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

var decoder = schema.NewDecoder()
var router = httprouter.New()

func Start(addr string) {
	decoder.SetAliasTag("json")
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
