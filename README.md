Embeddable swagger-ui for your Go program
=========================================

[![Go Reference](https://pkg.go.dev/badge/github.com/mwmahlberg/swagger-ui.svg)][godoc]
[![Quality Gate Status][qualitygate]][swaggerui:sonarcloud]

This [Go][go] module allows you to conveniently serve a swagger UI embedded into
your server.

Example
-------

```go [main.go]

package main

import (
  _ "embed"
  "flag"
  "log"
  "net/http"

  swaggerui "github.com/mwmahlberg/swagger-ui"
)

//go:embed petstore.yaml
var petStore []byte

var bind string

func init() {
  flag.StringVar(&bind, "bind", ":8080", "address to bind to")
}

func main() {

  flag.Parse()

  ui, err := swaggerui.New(swaggerui.Spec(swaggerui.DefaultSpecfileName, petStore))
  if err != nil {
    log.Fatalln(err)
  }

  mux := http.NewServeMux()
  mux.Handle("/api-docs/", http.StripPrefix("/api-docs/", ui))

  mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("dummy api"))
  })

  http.ListenAndServe(bind, mux)
}


```

Links
-----

- [Go programming language][go]
- [SwaggerUI project](https://github.com/swagger-api/swagger-ui)

[go]: https://go.dev "Project page of the Go programming language"
[godoc]: https://pkg.go.dev/github.com/mwmahlberg/swagger-ui
[qualitygate]: https://sonarcloud.io/api/project_badges/measure?project=mwmahlberg_swagger-ui&metric=alert_status
[swaggerui:sonarcloud]: https://sonarcloud.io/summary/new_code?id=mwmahlberg_swagger-ui
