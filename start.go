package idek

import (
	"log"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"
)

var decoder = schema.NewDecoder()
var router = httprouter.New()

func Start(addr string) {
	decoder.SetAliasTag("json")
	decoder.IgnoreUnknownKeys(true)

	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
