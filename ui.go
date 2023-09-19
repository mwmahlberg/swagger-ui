/*
 *  ui.go is part of github.com/mwmahlberg/swagger-ui project.
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

/*
 *  ui.go is part of github.com/mwmahlberg/swagger-ui project.
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

package swaggerui

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/asaskevich/govalidator"
	"github.com/mwmahlberg/memfs"
	"github.com/yalue/merged_fs"
)

const (
	// DefaultSpecfileName is the default name of the spec file.
	DefaultSpecfileName string = "swagger.yaml"
	// InitializerFilename is the default name of the initializer file.
	InitializerFilename string = "swagger-initializer.js"

	// This is the template for the initializer.js file, which is used to initialize the swagger-ui.
	// Alternatively, you can provide your own initializer by using the InitializerContent option.
	InitializerTemplate string = `
window.onload = function () {
  //<editor-fold desc="Changeable Configuration Block">

  // the following lines will be replaced by docker/configurator, when it runs in a docker-container
  window.ui = SwaggerUIBundle({
    url: "{{- if .Prefix -}}{{.Prefix}}{{- else -}}.{{- end -}}/{{.Filename}}",
    dom_id: '#swagger-ui',
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });

  //</editor-fold>
};
`
	embedPrefix string = "swagger-ui/dist"
)

//go:embed swagger-ui/dist
var swaggerui embed.FS

// HandlerOptions are used to initialize the handler.
type HandlerOption func(*SwaggerUi)

// SwaggerUi is a handler that serves the swagger-ui.
type SwaggerUi struct {
	// specFilename string `valid:"acceptedFileName~File name is wrong"`

	Overlay    *memfs.FS           `valid:"-"` // The overlay fs that is used to serve the custom spec and initializer
	Static     *fs.FS              `valid:"-"` // The base fs that is used to serve the static files of swagger-ui
	Merged     *merged_fs.MergedFS `valid:"-"` // The overlayfs that is used to serve the swagger-ui
	fileServer http.Handler        `valid:"-"` // The fileserver that is used to serve the swagger-ui

	specFilename string `valid:"stringlength(1|255)~File name is wrong)"`
	specContent  []byte `valid:"correctContent~File content is wrong"`

	initializerContent []byte `valid:"length(249|16384)~Initializer too small"` // The min length is the length of a minified version of a swagger-initializer

}

// ServeHTTP implements the http.Handler interface.
// It serves the swagger-ui, the spec file and the initializer by using the merged fs via http.FileServer.
func (ui *SwaggerUi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ui.fileServer.ServeHTTP(w, r)
}

// New returns a new SwaggerUi handler.
func New(opts ...HandlerOption) (*SwaggerUi, error) {
	var ui = &SwaggerUi{
		specFilename: DefaultSpecfileName}

	for _, opt := range opts {
		opt(ui)
	}

	if len(ui.initializerContent) == 0 {
		ui.initializerContent = getInitializer(ui.specFilename, "")
	}

	if isValid, err := govalidator.ValidateStruct(ui); !isValid {
		fmt.Println(err)
		return nil, SetupError{Cause: errors.New("invalid options: " + err.Error())}
	}

	if err := ui.setupOverlay(); err != nil {
		return nil, SetupError{Cause: errors.New("error setting up overlay: " + err.Error())}
	}

	if err := ui.setupStatic(); err != nil {
		return nil, SetupError{Cause: errors.New("error setting up static: " + err.Error())}
	}
	ui.Merged = merged_fs.NewMergedFS(fs.FS(ui.Overlay), *ui.Static)
	ui.fileServer = http.FileServer(http.FS(ui.Merged))
	return ui, nil
}

// Sets the name under which the data will be served as a spec file.
func Spec(name string, data []byte) HandlerOption {
	return func(suh *SwaggerUi) {
		suh.specFilename = name
		suh.specContent = data
	}
}

// Instead of using the default "swagger-initializer.js", you can provide your own.
func InitializerContent(code []byte) HandlerOption {
	return func(suh *SwaggerUi) {
		suh.initializerContent = code
	}
}

// Returns the name of the spec file served by the handler.
func (ui *SwaggerUi) SpecFilename() string {
	return ui.specFilename
}

func (ui *SwaggerUi) setupOverlay() error {
	o := memfs.New()

	if err := o.WriteFile(ui.specFilename, ui.specContent, 0644); err != nil {
		return errors.New("error writing specfile: " + err.Error())
	}

	if err := o.WriteFile(InitializerFilename, ui.initializerContent, 0644); err != nil {
		return errors.New("error writing initializer: " + err.Error())
	}

	ui.Overlay = o
	return nil
}

func (ui *SwaggerUi) setupStatic() (err error) {
	sub, err := fs.Sub(swaggerui, embedPrefix)
	if err != nil {
		return errors.New("error setting up base: " + err.Error())
	}
	ui.Static = &sub
	return nil
}

func getInitializer(filename string, prefix string) []byte {
	tmpl, _ := template.New(InitializerFilename).Parse(InitializerTemplate)
	var rendered bytes.Buffer
	tmpl.Execute(&rendered, struct {
		Prefix   string
		Filename string
	}{Prefix: prefix, Filename: filename})

	return rendered.Bytes()
}
