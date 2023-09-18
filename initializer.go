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
