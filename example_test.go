package swaggerui_test

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"

	swaggerui "github.com/mwmahlberg/swagger-ui"
)

func ExampleSwaggerUi() {
	_, _ = swaggerui.New()
	fmt.Println("test")
	// Output: test
}

func ExampleNew() {
	_, _ = swaggerui.New()
	fmt.Println("test2")
	// Output: test2
}

func ExampleSpec() {
	_, _ = swaggerui.New(swaggerui.Spec("foo.yaml", []byte("bar")))
	fmt.Println("test2")
	// Output: test2
}

func ExampleSwaggerUi_FileServer() {
	ui, err := swaggerui.New(swaggerui.Spec("foo.yaml", []byte("bar")))
	if err != nil {
		panic(err)
	}

	fs := ui.FileServer()
	ts := httptest.NewServer(fs)
	defer ts.Close()
	resp, err := http.Get(ts.URL + "/foo.yaml")
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		log.Fatalf("response content: %s", string(b))
	}

	fmt.Println(resp.StatusCode)
	// Output: 200
}

func ExampleSwaggerUi_FileSystem() {
	ui, err := swaggerui.New(swaggerui.Spec("foo.yaml", []byte("bar")))
	if err != nil {
		panic(err)
	}
	fs.WalkDir(ui.FileSystem(), ".", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, "foo.yaml") {
			fmt.Println(path)
		}
		if strings.HasSuffix(path, "swagger-initializer.js") {
			fmt.Println(path)
		}
		return nil
	})
	// Output: foo.yaml
	// swagger-initializer.js
}