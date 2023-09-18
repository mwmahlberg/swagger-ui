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
	"bytes"
	"html/template"
)

// This is the template for the initializer.js file, which is used to initialize the swagger-ui.
var InitializerTemplate = `
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

func getInitializer(filename string, prefix string) []byte {
	tmpl, _ := template.New("swagger-initializer.js").Parse(InitializerTemplate)
	var rendered bytes.Buffer
	tmpl.Execute(&rendered, struct {
		Prefix   string
		Filename string
	}{Prefix: prefix, Filename: filename})

	return rendered.Bytes()
}
