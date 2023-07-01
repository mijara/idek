package main

import (
	"idek"
	"idek/cmd/middleware"
	"idek/cmd/view"
	"idek/cmd/xerrors"
)

func main() {
	idek.Middleware(middleware.ValidateHeaders)
	idek.Error(xerrors.HandleError)

	idek.View("GET", "/hello/:name", view.Hello)
	idek.View("GET", "/rand", view.Rand)

	idek.ListenAndServe(":8080")
}
