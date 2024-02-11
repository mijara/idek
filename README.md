# Idek

A very opinionated web server framework with Go Generics at its core.

```go
package main

import (
	"fmt"
	"idek"
)

type HelloInput struct {
	Name string `json:"name"`
}

type HelloOutput struct {
	Message string `json:"message"`
}

func Hello(ctx *idek.Context, input HelloInput) (HelloOutput, error) {
	return HelloOutput{
		Message: fmt.Sprintf("Hello, %s!", input.Name),
	}, nil
}

func main() {
	idek.ViewHandler("GET", "/hello/:name", Hello)
	idek.Start(":8080")
}
```

Check `/cmd` for a more complete example.

## Opinions

- URL Params, Query params and body will all parse into the same object, you can pass them in any way you want. 

## Idek?

I don't even know...
