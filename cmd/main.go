package main

import (
	"idek"
	"idek/cmd/middleware"
	"idek/cmd/view"
	"idek/cmd/xerrors"
	"net/http"
)

func main() {
	idek.Middleware(middleware.ValidateHeaders)
	idek.Error(xerrors.HandleError)

	idek.View("/hello", view.Hello)
	idek.View("/rand", view.Rand)

	http.ListenAndServe(":8080", nil)
}
