/* 
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

package swaggerui

import (
	"embed"
	"errors"
	"io/fs"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/jncornett/overlayfs"
	"github.com/mwmahlberg/memfs"
)

//go:embed swagger-ui/dist
var swaggerui embed.FS

// HandlerOptions are used to initialize the handler.
type HandlerOption func(*SwaggerUi)

const (
	DefaultSpecfileName = "swagger.yaml"
	InitializerFilename = "swagger-initializer.js"
)

var (
	ErrEmptyData     = errors.New("no data was provided for the spec file")
	ErrInvalidJson   = errors.New("input is not valid JSON data")
	ErrEmptyFileName = errors.New("empty file name was provided")
	ErrNoJsonSuffix  = errors.New("name JSON spec does not end in '.json'")
)

// SwaggerUi is a handler that serves the swagger-ui.
type SwaggerUi struct {
	specFilename       string `valid:"matchingFilename"`
	specContent        []byte `valid:"correctContent"`
	initializerContent []byte `valid:"length(249)"` // The min length is the length of a minified version of a swagger-initializer
	fs                 fs.FS  `valid:"-"`
}

// New returns a new SwaggerUi handler.
func New(opts ...HandlerOption) (*SwaggerUi, error) {
	var ui = &SwaggerUi{
		specFilename: DefaultSpecfileName}

	for _, opt := range opts {
		opt(ui)
	}

	if len(ui.initializerContent) == 0 {
		if ui.specFilename == "" {
			return nil, SetupError{Cause: errors.New("no specfilename and no initializer given")}
		}
		ui.initializerContent = getInitializer(ui.specFilename, "")
	}

	if isValid, err := govalidator.ValidateStruct(ui); !isValid {
		return nil, SetupError{Cause: err}
	}
	return ui, nil
}

// Sets the name under which the data will be served as a spec file.
func Spec(name string, data []byte) HandlerOption {
	return func(suh *SwaggerUi) {
		suh.specFilename = name
		suh.specContent = data
	}
}

// Instead of using InitializerTemplate, you can provide your own.
func InitializerContent(code []byte) HandlerOption {
	return func(suh *SwaggerUi) {
		suh.initializerContent = code
	}
}

// FileSystem returns the http.FileSystem that is used to serve the swagger-ui.
func (ui *SwaggerUi) FileSystem() fs.FS {

	// Ensure this is a singleton
	if ui.fs == nil {
		ui.fs = ui.setupfs()
	}
	return ui.fs
}

// FileServer returns a http.Handler that serves the swagger-ui.
func (ui *SwaggerUi) FileServer() http.Handler {
	strippedUi, _ := fs.Sub(swaggerui, "swagger-ui/dist")
	ofs := overlayfs.NewOverlayFs(http.FS(strippedUi), http.FS(ui.FileSystem()))
	return http.FileServer(ofs)
}

func (ui *SwaggerUi) setupfs() fs.FS {
	overlay := memfs.New()

	overlay.WriteFile(ui.specFilename, ui.specContent, 0644)

	overlay.WriteFile("swagger-initializer.js", ui.initializerContent, 0644)

	return overlay
}
