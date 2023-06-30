package idek

import "github.com/gorilla/schema"

var decoder = schema.NewDecoder()

func init() {
	decoder.SetAliasTag("json")
}
