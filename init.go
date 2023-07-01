package idek

import (
	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

var decoder = schema.NewDecoder()
var router = httprouter.New()

func init() {
	decoder.SetAliasTag("json")
}

func ListenAndServe(addr string) {
	http.ListenAndServe(addr, router)
}
