package idek

import (
	"flag"
	"github.com/gorilla/schema"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

var decoder = schema.NewDecoder()
var router = httprouter.New()

func Start(addr string) {
	docsFlag := flag.Bool("docs", false, "generate docs")
	flag.Parse()

	if *docsFlag {
		data := generateDocsJSON()
		if err := os.WriteFile("docs.json", data, 0600); err != nil {
			log.Fatal(err)
		}
	} else {
		decoder.SetAliasTag("json")
		if err := http.ListenAndServe(addr, router); err != nil {
			log.Fatal(err)
		}
	}
}
