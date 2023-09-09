package swaggerui

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strconv"
)

//go:embed swagger-ui/dist
var swaggerui embed.FS
var ErrInvalidJson = errors.New("input is not valid JSON data")

func Handler(swaggerfile []byte) (http.Handler, http.FileSystem, error) {
	if !json.Valid(swaggerfile) {
		return nil, nil, ErrInvalidJson
	}

	fs.WalkDir(swaggerui, ".", func(path string, d fs.DirEntry, err error) error {
		log.Println(path)
		return nil
	})

	file := swaggerHandler(swaggerfile)
	ui, err := swaggerUIhandler()
	if ui == nil && err != nil {
		return nil, nil, fmt.Errorf("creating ui: %s", err)
	}

	return file, ui, nil
}

func swaggerUIhandler() (http.FileSystem, error) {
	fsys, err := fs.Sub(swaggerui, "swagger-ui/dist")
	if err != nil {
		return nil, err
	}
	fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		log.Println(path)
		return nil
	})
	return http.FS(fsys), nil
}

func swaggerHandler(swaggerfile []byte) http.Handler {
	contentLength := strconv.Itoa(len(swaggerfile))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", contentLength)
		w.Header().Set("Content-Type", "application/json")
		w.Write(swaggerfile)
	})
}
