/*
 *  example_test.go is part of github.com/mwmahlberg/swagger-ui project.
 *
 *  Copyright 2023 Markus W Mahlberg
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package swaggerui_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	swaggerui "github.com/mwmahlberg/swagger-ui"
)

const (
	someYaml string = `
openapi: 3.0.0
info:
	title: Swagger Petstore
	description: This is a sample server Petstore server.
	termsOfService: http://swagger.io/terms/
`
)

func ExampleSwaggerUi() {
	ui, err := swaggerui.New(swaggerui.Spec(swaggerui.DefaultSpecfileName, []byte(someYaml)))
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/api-docs/", http.StripPrefix("/api-docs/", ui))
	ts := httptest.NewServer(mux)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/api-docs/swagger.yaml")
	if err != nil {
		panic(err)
	}
	fmt.Println(resp.StatusCode)
	// Output: 200

}

// Use the Spec option to provide a swagger spec.
func ExampleSpec() {
	// Error handling ommitted for brevity

	// This will serve the spec provided under the default filename "swagger.yaml"
	ui, _ := swaggerui.New(swaggerui.Spec(swaggerui.DefaultSpecfileName, []byte("bar")))

	// Set up a mux and use the handler to serve the UI under /api-docs/
	mux := http.NewServeMux()
	mux.Handle("/api-docs/", http.StripPrefix("/api-docs/", ui))

	// Start a test server
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Get the spec
	resp, _ := http.Get(ts.URL + "/api-docs/" + swaggerui.DefaultSpecfileName)
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Printf("%s: %s", ui.SpecFilename(), string(body))
	// Output: swagger.yaml: bar
}

func ExampleInitializerContent() {
	// Error handling ommitted for brevity

	// This will serve an empty initializer
	// Also, the spec file will be named "swagger.yaml" as per the default
	ui, _ := swaggerui.New(swaggerui.InitializerContent([]byte("{}")))

	// Set up a mux and use the handler to serve the UI under /api-docs/
	mux := http.NewServeMux()
	mux.Handle("/api-docs/", http.StripPrefix("/api-docs/", ui))

	// Start a test server
	ts := httptest.NewServer(mux)
	defer ts.Close()

	// Get the initializer
	resp, _ := http.Get(ts.URL + "/api-docs/swagger-initializer.js")

	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	fmt.Printf("%s", string(body))
	// Output: {}
}

// func ExampleSwaggerUi_FileServer() {
// 	ui, err := swaggerui.New(swaggerui.Spec("foo.yaml", []byte("bar")))
// 	if err != nil {
// 		panic(err)
// 	}

// 	fs := ui.FileServer()
// 	ts := httptest.NewServer(fs)
// 	defer ts.Close()
// 	resp, err := http.Get(ts.URL + "/foo.yaml")
// 	if err != nil {
// 		panic(err)
// 	}

// 	if resp.StatusCode != 200 {
// 		b, _ := io.ReadAll(resp.Body)
// 		log.Fatalf("response content: %s", string(b))
// 	}

// 	fmt.Println(resp.StatusCode)
// 	// Output: 200
// }

// func ExampleSwaggerUi_FileSystem() {
// 	ui, err := swaggerui.New(swaggerui.Spec("foo.yaml", []byte("bar")))
// 	if err != nil {
// 		panic(err)
// 	}
// 	fs.WalkDir(ui.FileSystem(), ".", func(path string, d fs.DirEntry, err error) error {
// 		if strings.HasSuffix(path, "foo.yaml") {
// 			fmt.Println(path)
// 		}
// 		if strings.HasSuffix(path, "swagger-initializer.js") {
// 			fmt.Println(path)
// 		}
// 		return nil
// 	})
// 	// Output: foo.yaml
// 	// swagger-initializer.js
// }
